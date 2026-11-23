package events

// Package events defines the event contracts shared between services.
// Treat these like a public API: evolve them in a backward-compatible way.

import "time"

const (
	// TopicOrderCreated is the topic where the Order service publishes facts.
	TopicSatelitesCreated = "satelites.created"
	TopicRocketsCreated   = "rockets.created"
	// Consumer groups = "subscriptions". Each group gets its own copy of every event.
	GroupHyperSatelite = "hyper-satelite-service"
	GroupHyperRocket   = "hyper-rocket-service"
)

type Rocket struct {
	Satelite int `json:"satelite"`
	Id       int `json:"id"`
}

type Satelite struct {
	Tenant int `json:"tenant"`
	Id     int `json:"id"`
}

// OrderCreated is an immutable fact: "this order was created".
type RocketCreated struct {
	EventID   string    `json:"event_id"` // lets consumers deduplicate (idempotency)
	Rocket_pa *Rocket   `json:"rocket_p"`
	CreatedAt time.Time `json:"created_at"`
}

type SateliteCreated struct {
	EventID     string    `json:"event_id"` // lets consumers deduplicate (idempotency)
	Satelite_pa *Satelite `json:"satelite_p"`
	CreatedAt   time.Time `json:"created_at"`
}
