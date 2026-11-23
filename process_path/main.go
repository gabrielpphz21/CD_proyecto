package main

import (
	"encoding/json"
	"math/rand"
	"time"
)

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

type HyperSatelite struct {
	Id     int     `json:"id"`
	Energy int     `json:"energy"`
	Parts  string  `json:"parts"`
	State  string  `json:"state"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Z      float64 `json:"z"`
	Date   string  `json:"date"`
}

var state_values = []string{"OPTIMAL", "DEAD", "OFFLINE", "DAMAGED"}

func RandomDate(start, end time.Time) string {
	days := int(end.Sub(start).Hours() / 24)
	randomDays := rand.Intn(days + 1)

	return start.AddDate(0, 0, randomDays).Format("2006-01-02")
}

func MakeHyperSatelite(id int) (*HyperSatelite, error) {
	var p HyperSatelite
	p.Id = id
	p.X = 0.07 + rand.Float64()*(0.465-0.07)
	p.Y = 0.07 + rand.Float64()*(0.465-0.07)
	p.Z = 0.07 + rand.Float64()*(0.465-0.07)
	p.Energy = rand.Intn(100)
	p.Date = RandomDate(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC))

	m := map[string]string{
		"solar cells":    "healthy",
		"computer":       "healthy",
		"energy system":  "healthy",
		"cooling system": "healthy",
	}

	data, err := json.Marshal(m)

	if err != nil {
		return nil, err
	}

	p.Parts = string(data)
	return &p, nil
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
