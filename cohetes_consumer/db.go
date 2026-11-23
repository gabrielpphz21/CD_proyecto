package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

func conn() (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), ("postgres://postgres:password@psqlDB:5432/dyson_rings"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfull db connection")
	return conn, err
}

func insert_hyperrocket(h_rocket *HyperRocket, tx pgx.Tx) error {
	query := "INSERT INTO hyper_rocket (id, fuel, state, parts, x, y, z, date) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"

	parts_state := json.RawMessage(h_rocket.Parts)
	_, err := tx.Exec(context.Background(), query, h_rocket.Id, h_rocket.Fuel,
		h_rocket.State, parts_state, h_rocket.X, h_rocket.Y, h_rocket.Z, h_rocket.Date)
	return err
}
