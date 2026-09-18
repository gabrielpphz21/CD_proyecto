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

func heartbeatHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("backend alive"))
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

	mux.HandleFunc("/heartbeat", withInstanceHeader(heartbeatHandler))

	err1 := http.ListenAndServe(":8083", mux)
	if err1 != nil {
		log.Fatal(err1)
	}

}
