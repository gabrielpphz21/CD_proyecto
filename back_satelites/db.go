package main

import (
	"context"
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
