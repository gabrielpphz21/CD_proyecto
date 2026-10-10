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

	mux := http.NewServeMux()
	conn, err := conn()
	if err != nil {
		log.Fatal("Error al conectar con la BD")
	}

	mux.HandleFunc("/rocket", withInstanceHeader(func(w http.ResponseWriter, r *http.Request) {
		rocket_Handler(conn, w, r)
	}))

	mux.HandleFunc("/rocket/", withInstanceHeader(func(w http.ResponseWriter, r *http.Request) {
		rocket_Handler(conn, w, r)
	}))

	mux.HandleFunc("/hyper-rocket", withInstanceHeader(func(w http.ResponseWriter, r *http.Request) {
		HyperRocketHandler(conn, w, r)
	}))

	mux.HandleFunc("/hyper-rocket/", withInstanceHeader(func(w http.ResponseWriter, r *http.Request) {
		HyperRocketHandler(conn, w, r)
	}))

	log.Fatal(http.ListenAndServe(":8000", mux))

}
