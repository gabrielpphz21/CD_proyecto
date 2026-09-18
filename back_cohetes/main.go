package main

import (
	"log"
	"net/http"
)

func main() {

	mux := http.NewServeMux()
	conn, err := conn()
	if err != nil {
		log.Fatal("Error al conectar con la BD")
	}

	mux.HandleFunc("/rocket", func(w http.ResponseWriter, r *http.Request) {
		rocket_Handler(conn, w, r)
	})

	mux.HandleFunc("/rocket/", func(w http.ResponseWriter, r *http.Request) {
		rocket_Handler(conn, w, r)
	})

	mux.HandleFunc("/hyper-rocket", func(w http.ResponseWriter, r *http.Request) {
		HyperRocketHandler(conn, w, r)
	})
	mux.HandleFunc("/hyper-rocket/", func(w http.ResponseWriter, r *http.Request) {
		HyperRocketHandler(conn, w, r)
	})

	err2 := http.ListenAndServe(":8081", mux)
	if err2 != nil {
		log.Fatal(err2)
	}

}
