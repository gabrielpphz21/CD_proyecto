package main

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5"
)

type Rocket struct {
	Satelite int `json:"satelite"`
	Id       int `json:"id"`
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

func insert_rocket(rocket *Rocket, connection *pgx.Conn) error {
	query := "INSERT INTO rockets (satelite) VALUES ($1)"
	_, err := connection.Exec(context.Background(), query, rocket.Satelite)
	return err
}

func delete_rocket(id string, connection *pgx.Conn) error {
	query := "DELETE FROM rockets WHERE id = $1"
	_, err := connection.Exec(context.Background(), query, id)
	return err
}

func update_rocket(rocket *Rocket, connection *pgx.Conn) error {
	query := "UPDATE rockets SET satelite =  $1 WHERE id = $2"
	r, err := connection.Exec(context.Background(), query, rocket.Satelite, rocket.Id)
	if err != nil {
		return err
	}
	if r.RowsAffected() == 0 {
		return errors.New("The rocket whith provided id does not exist")
	}
	return nil

}

func get_rockets(connection *pgx.Conn) ([]Rocket, error) {
	query := "SELECT * FROM rockets LIMIT 100"
	rows, err := connection.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rockets []Rocket
	for rows.Next() {
		var rocket Rocket

		err1 := rows.Scan(
			&rocket.Id,
			&rocket.Satelite,
		)

		if err1 != nil {
			return nil, err1
		}

		rockets = append(rockets, rocket)
	}

	return rockets, nil
}

func insert_hyperrocket(h_rocket *HyperRocket, connection *pgx.Conn) error {
	query := "INSERT INTO hyper_rocket (id, fuel, state, parts, x, y, z, date) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"

	parts_state := json.RawMessage(h_rocket.Parts)
	_, err := connection.Exec(context.Background(), query, h_rocket.Id, h_rocket.Fuel,
		h_rocket.State, parts_state, h_rocket.X, h_rocket.Y, h_rocket.Z, h_rocket.Date)
	return err
}

func delete_hyperrocketLatest(id string, connection *pgx.Conn) error {
	query := "DELETE FROM hyper_rocket WHERE id = $1 AND date= (SELECT MAX (date) FROM hyper_rocket WHERE id=$1)"

	_, err := connection.Exec(context.Background(), query, id)
	return err
}

func rocket_info(id string, connection *pgx.Conn) (*Rocket, []HyperRocket, error) {
	query := "SELECT * FROM hyper_rocket WHERE id=$1 ORDER BY date desc"
	var h_rockets []HyperRocket
	rows, err := connection.Query(context.Background(), query, id)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var h_rocket HyperRocket

		err1 := rows.Scan(&h_rocket.Id, &h_rocket.Fuel, &h_rocket.Parts, &h_rocket.State, &h_rocket.X, &h_rocket.Y,
			&h_rocket.Z, &h_rocket.Date,
		)

		if err1 != nil {
			return nil, nil, err1
		}

		h_rockets = append(h_rockets, h_rocket)
	}

	var rocket Rocket
	query2 := "SELECT * FROM rockets WHERE id=$1 LIMIT 1"
	rows2, err2 := connection.Query(context.Background(), query2, id)
	if err2 != nil {
		return nil, nil, err2
	}
	defer rows2.Close()
	id_exists := false
	for rows2.Next() {

		err3 := rows2.Scan(&rocket.Id, &rocket.Satelite)

		if err3 != nil {
			return nil, nil, err3
		}
		id_exists = true
	}
	if id_exists == false {
		return nil, nil, errors.New("Rocket whit provided id does not exist")
	}

	return &rocket, h_rockets, nil

}
