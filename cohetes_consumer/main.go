package main

import (
	"context"
	"encoding/json"
	"events"
	"fmt"
	"kafkax"
	"log"
	"math/rand"
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

func main() {
	conn_db, err_1 := conn()

	if err_1 != nil {
		log.Fatal("Error while initializing database")
	}

	broker := getenv("KAFKA_BROKER", "localhost:9092")
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log.Println("HyperRocket-service started")
	if err := kafkax.ConsumeDb(ctx, []string{broker}, events.TopicRocketsCreated, events.GroupHyperRocket, handle, conn_db); err != nil {
		log.Fatal(err)
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
