package main

import (
	"bytes"
	"context"
	"encoding/json"
	"events"
	"fmt"
	"kafkax"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/segmentio/kafka-go"
)

// In-memory dedup for the demo. In production use a DB unique constraint
// or an "processed_events" table.
var processed = map[string]struct{}{}

func handle(_ context.Context, msg kafka.Message, conn *pgx.Conn) error {
	var evt events.RocketCreated
	if err := json.Unmarshal(msg.Value, &evt); err != nil {
		return nil // poison message: skip (or send to DLQ)
	}

	if _, seen := processed[evt.EventID]; seen {
		log.Printf("[inventory] duplicate event %s ignored", evt.EventID)
		return nil
	}

	r_id := evt.Rocket_pa.Id
	var hr_arr = []HyperRocket{}

	for i := 0; i < 20; i++ {
		hr_i, err_i := MakeHyperRocket(r_id)
		if err_i != nil {
			fmt.Print("Error while making the hypersatelite info in kafka consumer")
			return nil
		}
		hr_arr = append(hr_arr, *hr_i)
	}

	tx, err1 := conn.Begin(context.Background())
	defer tx.Rollback(context.Background())
	if err1 != nil {
		return nil
	}

	for i := 0; i < 20; i++ {
		err_i2 := insert_hyperrocket(&hr_arr[i], tx)
		if err_i2 != nil {
			fmt.Print("Error while inserting the hypersatelite info in kafka consumer")
			return nil
		}
		log.Print("Insertando hyperrocket desde consumidor : ", hr_arr[i])

	}
	tx.Commit(context.Background())

	processed[evt.EventID] = struct{}{}
	return nil
}

type HyperRocket struct {
	Id    int     `json:"id"`
	Fuel  float64 `json:"fuel"`
	Parts string  `json:"parts"`
	State string  `json:"state"`
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Z     float64 `json:"z"`
	Date  string  `json:"date"`
}

func RandomDate(start, end time.Time) string {
	days := int(end.Sub(start).Hours() / 24)
	randomDays := rand.Intn(days + 1)

	return start.AddDate(0, 0, randomDays).Format("2006-01-02")
}

func MakeHyperRocket(id int) (*HyperRocket, error) {
	var p HyperRocket
	p.Id = id
	p.X = 0.07 + rand.Float64()*(0.465-0.07)
	p.Y = 0.07 + rand.Float64()*(0.465-0.07)
	p.Z = 0.07 + rand.Float64()*(0.465-0.07)
	p.Fuel = rand.Float64() * 100
	p.Date = RandomDate(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC))

	m := map[string]string{
		"engine":          "healthy",
		"computer":        "healthy",
		"fuelling system": "healthy",
		"cooling system":  "healthy",
	}

	data, err := json.Marshal(m)

	if err != nil {
		return nil, err
	}

	p.Parts = string(data)
	return &p, nil
}

var instanceID = "unknown"

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func registerWithConsul(consulAddr string) error {
	payload := map[string]any{
		"ID":      instanceID,
		"Name":    "cohetes-consumer",
		"Address": instanceID, // resolvable via Docker's embedded DNS
		"Port":    8080,
		"Check": map[string]any{
			"HTTP":                           fmt.Sprintf("http://%s:8080/health", instanceID),
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
	conn_db, err_1 := conn()
	consulAddr := getEnv("CONSUL_HTTP_ADDR", "http://consul:8500")
	instanceID = getEnv("INSTANCE_ID", "cohetes-consumer")

	// Start a lightweight health endpoint in the background (GPTOSO health check)
	go func() {
		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			// Optional: verify Kafka connection / consumer loop is healthy
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
		})
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Printf("health server error: %v", err)
		}
	}()

	go func() {
		if err := registerWithConsul(consulAddr); err != nil {
			log.Printf("WARNING: %v", err)
		}
	}()

	if err_1 != nil {
		log.Fatal("Error while initializing database")
	}

	broker := getEnv("KAFKA_BROKER", "localhost:9092")
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log.Println("HyperRocket-service started")
	if err := kafkax.ConsumeDb(ctx, []string{broker}, events.TopicRocketsCreated, events.GroupHyperRocket, handle, conn_db); err != nil {
		log.Fatal(err)
	}

	// Graceful shutdown: on SIGTERM/SIGINT (what `docker stop` sends),
	// deregister from Consul BEFORE the process actually exits, so
	// gateways stop getting routed to a backend that's already gone.
	log.Printf("shutting down %s", instanceID)
	deregisterFromConsul(consulAddr)

}
