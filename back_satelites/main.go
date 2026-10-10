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
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/satelite", withInstanceHeader(func(w http.ResponseWriter, r *http.Request) {
		SateliteHandler(conn, w, r)
	}))

	mux.HandleFunc("/satelite/", withInstanceHeader(func(w http.ResponseWriter, r *http.Request) {
		SateliteHandler(conn, w, r)
	}))

	mux.HandleFunc("/hyper-satelite", withInstanceHeader(func(w http.ResponseWriter, r *http.Request) {
		HyperSateliteHandler(conn, w, r)
	}))

	mux.HandleFunc("/hyper-satelite/", withInstanceHeader(func(w http.ResponseWriter, r *http.Request) {
		HyperSateliteHandler(conn, w, r)
	}))

	log.Fatal(http.ListenAndServe(":8000", mux))

}
