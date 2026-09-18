package main

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type EnergyTransaction struct {
	ID          int     `json:"id"`
	Volume      float64 `json:"volume"`
	Origin      int     `json:"origin"`
	Destination int     `json:"destination"`
	State       string  `json:"state"`
	Date        string  `json:"date"`
}

type Tenant struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Tax_id string `json:"tax_id"`
}

type Destination struct {
	Id     int    `json:"id"`
	Planet string `json:"planet"`
	State  string `json:"state"`
}

type TakeOff struct {
	Id     int    `json:"id"`
	Rocket int    `json:"rocket"`
	Date   string `json:"date"`
	Runway int    `json:"runway"`
}

func insert_tenant(tenant *Tenant, connection *pgx.Conn) error {
	query := "INSERT INTO tenants (name, tax_id) VALUES ($1, $2)"
	_, err := connection.Exec(context.Background(), query, tenant.Name, tenant.Tax_id)

	return err
}

func delete_tenant(tenant_id string, connection *pgx.Conn) error {
	query := "DELETE FROM tenants WHERE id = $1"
	_, err := connection.Exec(context.Background(), query, tenant_id)

	return err
}

func update_tenant(tenant *Tenant, connection *pgx.Conn) error {
	query := "UPDATE tenants SET name = $1, tax_id=$2 WHERE id = $3"
	r, err := connection.Exec(context.Background(), query, tenant.Name, tenant.Tax_id, tenant.Id)
	if r.RowsAffected() == 0 {
		return errors.New("Tenant with provided id does not exist")
	}

	return err
}

func get_Tenants(connection *pgx.Conn) ([]Tenant, error) {
	query := "SELECT * FROM tenants ORDER BY id DESC LIMIT 100"
	rows, err := connection.Query(context.Background(), query)
	var tenants []Tenant

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var tnt Tenant
		err1 := rows.Scan(&tnt.Id, &tnt.Name, &tnt.Tax_id)
		tenants = append(tenants, tnt)

		if err1 != nil {
			return nil, err1
		}
	}

	return tenants, err

}

func get_TenantId(connection *pgx.Conn, tenant_id string) (*Tenant, error) {
	query := "SELECT * FROM tenants WHERE id=$1 LIMIT 1"
	tnt := &Tenant{}
	err := connection.QueryRow(context.Background(), query, tenant_id).Scan(&tnt.Id, &tnt.Name, &tnt.Tax_id)
	if err != nil {
		return nil, err
	}
	return tnt, err

}

func insert_energyTransaction(e_transaction *EnergyTransaction, connection *pgx.Conn) error {
	query := "INSERT INTO energy_transactions (volume, origin, destination, state, date) VALUES ($1, $2, $3, $4, $5)"
	_, err := connection.Exec(context.Background(), query, e_transaction.Volume, e_transaction.Origin, e_transaction.Destination,
		e_transaction.State, e_transaction.Date)

	return err
}

func delete_energyTransaction(e_transaction_id string, connection *pgx.Conn) error {
	query := "DELETE FROM envergy_transactions WHERE id = $1"
	_, err := connection.Exec(context.Background(), query, e_transaction_id)

	return err
}

func update_energyTransaction(e_transaction *EnergyTransaction, connection *pgx.Conn) error {
	query := "UPDATE energy_transactions SET volume = $1, origin=$2, destination=$3, state=$4, date=$5 WHERE id = $6"
	r, err := connection.Exec(context.Background(), query, e_transaction.Volume, e_transaction.Origin, e_transaction.Destination,
		e_transaction.State, e_transaction.Date, e_transaction.ID)
	if r.RowsAffected() == 0 {
		return errors.New("Energy_transaction with provided id does not exist")
	}

	return err
}

func get_energyTransactions(connection *pgx.Conn) ([]EnergyTransaction, error) {
	query := "SELECT * FROM energy_transactions ORDER BY date DESC LIMIT 100"
	rows, err := connection.Query(context.Background(), query)
	var e_transactions []EnergyTransaction
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var e_t EnergyTransaction
		err1 := rows.Scan(&e_t.ID, &e_t.Volume, &e_t.Origin, &e_t.Destination, &e_t.State, &e_t.Date)
		e_transactions = append(e_transactions, e_t)

		if err1 != nil {
			return nil, err1
		}
	}

	return e_transactions, err

}

func get_energyTransactionsId(connection *pgx.Conn, energy_transaction_id string) (*EnergyTransaction, error) {
	e_t := &EnergyTransaction{}

	query := "SELECT * FROM energy_transactions WHERE id=$1 ORDER BY date DESC LIMIT 100"
	err := connection.QueryRow(context.Background(), query, energy_transaction_id).Scan(&e_t.ID,
		&e_t.Volume, &e_t.Origin, &e_t.Destination, &e_t.State, &e_t.Date)

	if err != nil {
		return e_t, nil
	}

	return nil, err

}

func insert_destination(destination *Destination, connection *pgx.Conn) error {
	query := "INSERT INTO destinations (planet, state) VALUES ($1, $2)"
	_, err := connection.Exec(context.Background(), query, destination.Planet, destination.State)

	return err
}

func delete_destination(destination_id string, connection *pgx.Conn) error {
	query := "DELETE FROM destinations WHERE id = $1"
	_, err := connection.Exec(context.Background(), query, destination_id)

	return err
}

func update_destinantion(destination *Destination, connection *pgx.Conn) error {
	query := "UPDATE destinations SET planet = $1, state=$2 WHERE id = $3"
	r, err := connection.Exec(context.Background(), query, destination.Planet, destination.State, destination.Id)
	if r.RowsAffected() == 0 {
		return errors.New("Destination with provided id does not exist")
	}

	return err
}

func get_destinations(connection *pgx.Conn) ([]Destination, error) {
	query := "SELECT * FROM destinations ORDER BY id DESC LIMIT 100"
	rows, err := connection.Query(context.Background(), query)
	var destinations []Destination
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var dtn Destination
		err1 := rows.Scan(&dtn.Id, &dtn.Planet, &dtn.State)
		destinations = append(destinations, dtn)

		if err1 != nil {
			return nil, err1
		}
	}

	return destinations, err

}

func get_destinationsId(connection *pgx.Conn, destination_id string) (*Destination, error) {
	query := "SELECT * FROM destinations WHERE id=$1 LIMIT 1"
	dtn := &Destination{}
	err := connection.QueryRow(context.Background(), query, destination_id).Scan(&dtn.Id, &dtn.Planet, &dtn.State)
	if err != nil {
		return nil, err
	}

	return dtn, nil
}

func insert_takeoff(takeoff *TakeOff, connection *pgx.Conn) error {
	query := "INSERT INTO take_offs (rocket, date, runway) VALUES ($1, $2, $3)"
	_, err := connection.Exec(context.Background(), query, takeoff.Rocket, takeoff.Date, takeoff.Runway)

	return err
}

func delete_takeoff(takeoff_id string, connection *pgx.Conn) error {
	query := "DELETE FROM take_offs WHERE id = $1"
	_, err := connection.Exec(context.Background(), query, takeoff_id)

	return err
}

func update_takeoff(takeoff *TakeOff, connection *pgx.Conn) error {
	query := "UPDATE take_offs SET rocket = $1, date=$2, runway=$3, WHERE id = $4"
	r, err := connection.Exec(context.Background(), query, takeoff.Rocket, takeoff.Date, takeoff.Runway)
	if r.RowsAffected() == 0 {
		return errors.New("Takeoff with provided id does not exist")
	}

	return err
}

func get_takeoffs(connection *pgx.Conn) ([]TakeOff, error) {
	query := "SELECT * FROM take_offs ORDER BY id DESC LIMIT 100"
	rows, err := connection.Query(context.Background(), query)
	var take_offs []TakeOff
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var tko TakeOff
		err1 := rows.Scan(&tko.Id, &tko.Rocket, &tko.Date, &tko.Runway)
		take_offs = append(take_offs, tko)

		if err1 != nil {
			return nil, err1
		}
	}

	return take_offs, err

}

func get_takeoffId(connection *pgx.Conn, take_off_id string) (*TakeOff, error) {
	query := "SELECT * FROM take_offs WHERE id=$1 LIMIT 1"
	tko := &TakeOff{}
	err := connection.QueryRow(context.Background(), query, take_off_id).Scan(&tko.Id, &tko.Rocket, &tko.Date, &tko.Runway)
	if err != nil {
		return nil, err
	}

	return tko, nil

}
