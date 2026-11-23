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

func insert_hypersatelites(h_satelite *HyperSatelite, tx pgx.Tx) error {
	query := "INSERT INTO hyper_satelite (id, energy, parts, state, x, y, z, date) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"

	parts_state := json.RawMessage(h_satelite.Parts)

	log.Printf(
		"QUERY  conn=%p query=%s",
		conn,
		query,
	)
	_, err := tx.Exec(context.Background(), query, h_satelite.Id, h_satelite.Energy,

		parts_state, h_satelite.State, h_satelite.X, h_satelite.Y, h_satelite.Z, h_satelite.Date)
	return err
}
