package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "protoc/proto_gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

var instanceID = "not known"

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func registerWithConsul(consulAddr string) error {
	payload := map[string]any{
		"ID":      instanceID,
		"Name":    "back-info",
		"Address": instanceID, // resolvable via Docker's embedded DNS
		"Port":    9000,
		"Check": map[string]any{
			"GRPC":                           fmt.Sprintf("%s:9000", instanceID),
			"GRPCUseTLS":                     false,
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

	db_conn, err := conn()

	if err != nil {
		log.Fatal("Base de datos caída")
	}

	lis, err1 := net.Listen("tcp", ":9000")
	if err1 != nil {
		log.Fatalf("failed to listen: %v", err1)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterDestinationServiceServer(grpcServer, &DestinationServer{conn: db_conn})
	pb.RegisterEnergyTransactionServiceServer(grpcServer, &EnergyTransactionServer{conn: db_conn})
	pb.RegisterTakeOffServiceServer(grpcServer, &TakeOffServer{conn: db_conn})
	pb.RegisterTenantServiceServer(grpcServer, &TenantServer{conn: db_conn})

	hs := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, hs)
	hs.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	reflection.Register(grpcServer)

	consulAddr := getEnv("CONSUL_HTTP_ADDR", "http://consul:8500")
	// Graceful shutdown: on SIGTERM/SIGINT (what `docker stop` sends),
	// deregister from Consul BEFORE the process actually exits, so
	// gateways stop getting routed to a backend that's already gone.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigCh)
	go func() {
		<-sigCh
		log.Printf("shutting down %s", instanceID)
		hs.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)
		deregisterFromConsul(consulAddr)

		stopped := make(chan struct{})
		go func() {
			grpcServer.GracefulStop()
			close(stopped)
		}()

		select {
		case <-stopped:
		case <-time.After(5 * time.Second):
			log.Printf("graceful stop timed out, forcing")
			grpcServer.Stop()
		}
	}()

	go func() {
		if err := registerWithConsul(consulAddr); err != nil {
			log.Printf("WARNING: %v", err)
		}
	}()

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("grpc server error: %v", err)
	}

}
