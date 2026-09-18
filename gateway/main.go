package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type RouteBeTargets struct {
	Route   string   `json:"route"`
	Targets []string `json:"targets"`
}

type BackendT struct {
	Path    *url.URL
	Proxy   *httputil.ReverseProxy
	healthy atomic.Bool
}

type Routegroup struct {
	Route_p  string
	Backends []*BackendT
	Counter  uint64
}

func (rg *Routegroup) next() *BackendT {
	n := len(rg.Backends)
	for i := 0; i < n; i++ {
		idx := atomic.AddUint64(&rg.Counter, 1) % uint64(n)
		b_t := rg.Backends[idx]
		if b_t.healthy.Load() {
			return b_t
		}
	}
	return nil
}

// ----------GATEWAY
type Gateway struct {
	groups []*Routegroup
	signal chan struct{}
}

func loadRouteBeTargets(path string) ([]RouteBeTargets, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfgs []RouteBeTargets
	if err := json.Unmarshal(data, &cfgs); err != nil {
		return nil, err
	}
	return cfgs, nil

}

func NewGateway(cfgs []RouteBeTargets, maxConcurrent int) *Gateway {
	var groups []*Routegroup

	for _, cfg := range cfgs {
		rg := &Routegroup{Route_p: cfg.Route}
		for _, t := range cfg.Targets {
			target, err := url.Parse(t)
			if err != nil {
				log.Fatalf("invalid target url %s: %v", t, err)
			}
			b := &BackendT{
				Path:  target,
				Proxy: httputil.NewSingleHostReverseProxy(target),
			}
			b.healthy.Store(true)
			rg.Backends = append(rg.Backends, b)
		}
		groups = append(groups, rg)
	}
	return &Gateway{
		groups: groups,
		signal: make(chan struct{}, maxConcurrent),
	}
}

func (g *Gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.signal <- struct{}{}
	defer func() { <-g.signal }()

	for _, rg := range g.groups {
		if !strings.HasPrefix(r.URL.Path, rg.Route_p) {
			continue
		}
		backend := rg.next()
		if backend == nil {
			http.Error(w, "no healthy backend available", http.StatusServiceUnavailable)
			return
		}
		log.Printf("[%s] %s -> %s (prefix %q)", r.Method, r.URL.Path, backend.Path, rg.Route_p)
		backend.Proxy.ServeHTTP(w, r)
		return
	}
	http.NotFound(w, r)
}

func (g *Gateway) startHeartbeatMonitor(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			var wg sync.WaitGroup
			for _, rg := range g.groups {
				for _, b := range rg.Backends {
					wg.Add(1)
					go func(b *BackendT) {
						defer wg.Done()
						b.healthy.Store(checkHeartbeat(b.Path.String()))
					}(b)
				}
			}
			wg.Wait()
		}
	}()
}

func checkHeartbeat(target string) bool {
	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(target + "/heartbeat")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

// statusHandler reports every backend's last-known health, grouped by
// route prefix, so you can watch load balancing / failover happen live.
func (g *Gateway) statusHandler(w http.ResponseWriter, r *http.Request) {
	type backendStatus struct {
		Target  string `json:"target"`
		Healthy bool   `json:"healthy"`
	}
	out := make(map[string][]backendStatus)
	for _, rg := range g.groups {
		var list []backendStatus
		for _, b := range rg.Backends {
			list = append(list, backendStatus{Target: b.Path.String(), Healthy: b.healthy.Load()})
		}
		out[rg.Route_p] = list
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

func heartbeatHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("gateway alive"))
}

func main() {
	cfgs, err := loadRouteBeTargets("routes.json")

	if err != nil {
		log.Fatalf("failed to load the routes file routes.json: %v", err)
	}
	gw := NewGateway(cfgs, 7)
	gw.startHeartbeatMonitor(5 * time.Second)

	mux := http.NewServeMux()
	mux.HandleFunc("/heartbeat", heartbeatHandler)
	mux.HandleFunc("/status", gw.statusHandler)
	mux.Handle("/", gw)

	log.Println("Gateway listening on: 8000")
	log.Fatal(http.ListenAndServe(":8000", mux))

}
