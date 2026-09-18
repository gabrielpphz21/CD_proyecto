package main

import (
	"log"
	"net/http"
)

func main() {

	conn, err := conn()

	if err != nil {
		log.Fatal("Base de datos caída")
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/energy-transaction", func(w http.ResponseWriter, r *http.Request) {
		EnergyTransactionHandler(conn, w, r)

	})

	mux.HandleFunc("/energy-transaction/", func(w http.ResponseWriter, r *http.Request) {
		EnergyTransactionHandler(conn, w, r)

	})

	mux.HandleFunc("/destination", func(w http.ResponseWriter, r *http.Request) {
		DestinationHandler(conn, w, r)
	})

	mux.HandleFunc("/destination/", func(w http.ResponseWriter, r *http.Request) {
		DestinationHandler(conn, w, r)
	})

	mux.HandleFunc("/take-off", func(w http.ResponseWriter, r *http.Request) {
		TakeOffHandler(conn, w, r)
	})

	mux.HandleFunc("/take-off/", func(w http.ResponseWriter, r *http.Request) {
		TakeOffHandler(conn, w, r)
	})

	mux.HandleFunc("/tenant", func(w http.ResponseWriter, r *http.Request) {
		TenantHandler(conn, w, r)
	})

	mux.HandleFunc("/tenant/", func(w http.ResponseWriter, r *http.Request) {
		TenantHandler(conn, w, r)
	})

	err1 := http.ListenAndServe(":8083", mux)
	if err1 != nil {
		log.Fatal(err1)
	}

}
