package main

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"
)

type Satelite struct {
	Tenant int `json:"tenant"`
	Id     int `json:"id"`
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

func insert_satelite(satelite *Satelite, connection *pgx.Conn) error {
	query := "INSERT INTO satelites (tenant) VALUES ($1)"
	_, err := connection.Exec(context.Background(), query, satelite.Tenant)
	return err
}

func delete_satelite(id string, connection *pgx.Conn) error {
	query := "DELETE FROM satelites WHERE id = $1"
	_, err := connection.Exec(context.Background(), query, id)
	return err
}

func update_satelite(satelite *Satelite, connection *pgx.Conn) error {
	query := "UPDATE satelites SET tenant =  $1 WHERE id = $2"
	_, err := connection.Exec(context.Background(), query, satelite.Tenant, satelite.Id)
	return err
}

func get_satelites(connection *pgx.Conn) ([]Satelite, error) {
	query := "SELECT * FROM satelites LIMIT 100"
	rows, err := connection.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var satelites []Satelite
	for rows.Next() {
		var satelite Satelite

		err1 := rows.Scan(
			&satelite.Id,
			&satelite.Tenant,
		)

		if err1 != nil {
			return nil, err1
		}

		satelites = append(satelites, satelite)
	}

	return satelites, nil
}

func insert_hypersatelites(h_satelite *HyperSatelite, connection *pgx.Conn) error {
	query := "INSERT INTO hyper_satelite (id, energy, parts, state, x, y, z, date) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"

	parts_state := json.RawMessage(h_satelite.Parts)
	_, err := connection.Exec(context.Background(), query, h_satelite.Id, h_satelite.Energy,
		parts_state, h_satelite.State, h_satelite.X, h_satelite.Y, h_satelite.Z, h_satelite.Date)
	return err
}

func delete_hypersateliteLatest(id string, connection *pgx.Conn) error {
	query := "DELETE FROM hyper_satelite WHERE id = $1 AND date= (SELECT MAX (date) FROM hyper_satelite WHERE id=$1)"

	_, err := connection.Exec(context.Background(), query, id)
	return err
}

func satelite_info(id string, connection *pgx.Conn) (*Satelite, []HyperSatelite, error) {
	query := "SELECT * FROM hyper_satelite WHERE id=$1 ORDER BY date desc"
	var h_satelites []HyperSatelite
	rows, err := connection.Query(context.Background(), query, id)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var h_satelite HyperSatelite
		err1 := rows.Scan(&h_satelite.Id, &h_satelite.Energy, &h_satelite.Parts, &h_satelite.State, &h_satelite.X, &h_satelite.Y,
			&h_satelite.Z, &h_satelite.Date,
		)

		if err1 != nil {
			return nil, nil, err1
		}

		h_satelites = append(h_satelites, h_satelite)
	}

	var satelite Satelite
	query2 := "SELECT * FROM satelites WHERE id=$1"
	rows2, err2 := connection.Query(context.Background(), query2, id)
	if err2 != nil {
		return nil, nil, err2
	}
	defer rows2.Close()
	for rows2.Next() {

		err3 := rows2.Scan(&satelite.Id, &satelite.Tenant)

		if err3 != nil {
			return nil, nil, err3
		}

	}

	return &satelite, h_satelites, nil

}
