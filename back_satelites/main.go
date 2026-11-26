package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"events"
	"fmt"
	"kafkax"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	pb "protoc/proto_gen"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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
		"Name":    "back-satelite",
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

func loadClientTLSCredentials() (credentials.TransportCredentials, error) {
	pemServerCA, err := os.ReadFile("/etc/ssl/certs/grpc/ca.crt")
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add server CA certificate")
	}

	clientCert, err := tls.LoadX509KeyPair("/etc/ssl/certs/grpc/tls.crt", "/etc/ssl/certs/grpc/tls.key")
	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
		ServerName:   "grpc-server-service", // Matches SAN in Server Cert
		MinVersion:   tls.VersionTLS13,
	}

	return credentials.NewTLS(config), nil
}

func main() {

	broker := getEnv("KAFKA_BROKER", "localhost:9092")

	creds, err := loadClientTLSCredentials()
	if err != nil {
		log.Fatalf("Failed to load client TLS credentials: %v", err)
	}

	// Connect to Kubernetes internal DNS service address
	target := "grpc-server-service.default.svc.cluster.local:50051"

	log.Printf("Connecting to %s with mTLS...", target)
	conn_mtls, err := grpc.NewClient(target, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn_mtls.Close()
	if err := kafkax.EnsureTopic(broker, events.TopicRocketsCreated, 3); err != nil {
		log.Fatalf("ensure topic: %v", err)
	}

	producer := kafkax.NewWriter([]string{broker}, events.TopicSatelitesCreated)
	defer producer.Close()

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

	pb.RegisterSateliteServiceServer(grpcServer, &Server{conn: db_conn, producer: producer})

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

	if errGRPC := grpcServer.Serve(lis); errGRPC != nil {
		log.Fatalf("grpc server error: %v", errGRPC)
	}

}
