package main

import (
	"log"
	"net/http"
)

var instanceID = "not known"

func withInstanceHeader(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Instance-Id", instanceID)
		next(w, r)
	}
}

func main() {

	conn, err := conn()

	if err != nil {
		log.Fatal("Base de datos caída")
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/energy-transaction", withInstanceHeader(func(w http.ResponseWriter, r *http.Request) {
		EnergyTransactionHandler(conn, w, r)

	}))

	mux.HandleFunc("/energy-transaction/", withInstanceHeader(func(w http.ResponseWriter, r *http.Request) {
		EnergyTransactionHandler(conn, w, r)

	}))

	mux.HandleFunc("/destination", withInstanceHeader(func(w http.ResponseWriter, r *http.Request) {
		DestinationHandler(conn, w, r)
	}))

	mux.HandleFunc("/destination/", withInstanceHeader(func(w http.ResponseWriter, r *http.Request) {
		DestinationHandler(conn, w, r)
	}))

	mux.HandleFunc("/take-off", withInstanceHeader(func(w http.ResponseWriter, r *http.Request) {
		TakeOffHandler(conn, w, r)
	}))

	mux.HandleFunc("/take-off/", withInstanceHeader(func(w http.ResponseWriter, r *http.Request) {
		TakeOffHandler(conn, w, r)
	}))

	mux.HandleFunc("/tenant", withInstanceHeader(func(w http.ResponseWriter, r *http.Request) {
		TenantHandler(conn, w, r)
	}))

	mux.HandleFunc("/tenant/", withInstanceHeader(func(w http.ResponseWriter, r *http.Request) {
		TenantHandler(conn, w, r)
	}))

	log.Fatal(http.ListenAndServe(":8000", mux))

}
