package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var instanceID = "not known"

func withInstanceHeader(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Instance-Id", instanceID)
		next(w, r)
	}
}

func heartbeatHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("backend alive"))
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func registerWithConsul(consulAddr string) error {
	payload := map[string]any{
		"ID":      instanceID,
		"Name":    "back_satelite",
		"Address": instanceID, // resolvable via Docker's embedded DNS
		"Port":    8080,
		"Check": map[string]any{
			"HTTP":                           fmt.Sprintf("http://%s:8080/heartbeat", instanceID),
			"Interval":                       "5s",
			"Timeout":                        "2s",
			"DeregisterCriticalServiceAfter": "1m",
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	var lastErr error
	for attempt := 1; attempt <= 10; attempt++ {
		req, err := http.NewRequest(http.MethodPut, consulAddr+"/v1/agent/service/register", bytes.NewReader(body))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				log.Printf("registered %s with Consul at %s", instanceID, consulAddr)
				return nil
			}
			lastErr = fmt.Errorf("consul returned status %d", resp.StatusCode)
		} else {
			lastErr = err
		}
		time.Sleep(time.Duration(attempt) * time.Second)
	}
	return fmt.Errorf("giving up registering with consul: %w", lastErr)
}

// deregisterFromConsul is called on shutdown so Consul (and therefore
// every gateway watching it) stops routing traffic here immediately,
// instead of waiting for DeregisterCriticalServiceAfter to time it out.
func deregisterFromConsul(consulAddr string) {
	req, err := http.NewRequest(http.MethodPut, consulAddr+"/v1/agent/service/deregister/"+instanceID, nil)
	if err != nil {
		log.Printf("failed to build deregister request: %v", err)
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("failed to deregister %s from consul: %v", instanceID, err)
		return
	}
	resp.Body.Close()
	log.Printf("deregistered %s from consul", instanceID)
}

func main() {
	if v := os.Getenv("INSTANCE_ID"); v != "" {
		instanceID = v
	}
	consulAddr := getEnv("CONSUL_HTTP_ADDR", "http://consul:8500")

	if err := registerWithConsul(consulAddr); err != nil {
		// Non-fatal: the service still runs and serves requests, it just
		// won't be discoverable via Consul until this succeeds. Logged
		// loudly because silently running "invisible" is a confusing
		// failure mode.
		log.Printf("WARNING: %v", err)
	}

	conn, err := conn()
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/satelite", withInstanceHeader(func(w http.ResponseWriter, r *http.Request) {
		SateliteHandler(conn, w, r)
	}))

	mux.HandleFunc("/satelite/", withInstanceHeader(func(w http.ResponseWriter, r *http.Request) {
		SateliteHandler(conn, w, r)
	}))

	mux.HandleFunc("/hyper-satelite", withInstanceHeader(func(w http.ResponseWriter, r *http.Request) {
		HyperSateliteHandler(conn, w, r)
	}))

	mux.HandleFunc("/hyper-satelite/", withInstanceHeader(func(w http.ResponseWriter, r *http.Request) {
		HyperSateliteHandler(conn, w, r)
	}))

	mux.HandleFunc("/heartbeat", withInstanceHeader(heartbeatHandler))

	srv := &http.Server{Addr: ":8080", Handler: mux}

	// Graceful shutdown: on SIGTERM/SIGINT (what `docker stop` sends),
	// deregister from Consul BEFORE the process actually exits, so
	// gateways stop getting routed to a backend that's already gone.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Printf("shutting down %s", instanceID)
		deregisterFromConsul(consulAddr)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		srv.Shutdown(ctx)
	}()

	log.Printf("backend %s listening on :8080", instanceID)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
