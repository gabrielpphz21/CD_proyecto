package main

import (
	"log"
	"net/http"
)

func main() {

	conn, err := conn()
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/satelite", func(w http.ResponseWriter, r *http.Request) {
		SateliteHandler(conn, w, r)
	})

	mux.HandleFunc("/satelite/", func(w http.ResponseWriter, r *http.Request) {
		SateliteHandler(conn, w, r)
	})

	mux.HandleFunc("/hyper-satelite", func(w http.ResponseWriter, r *http.Request) {
		HyperSateliteHandler(conn, w, r)
	})

	mux.HandleFunc("/hyper-satelite/", func(w http.ResponseWriter, r *http.Request) {
		HyperSateliteHandler(conn, w, r)
	})

	err1 := http.ListenAndServe(":8082", mux)
	if err1 != nil {
		log.Fatal(err1)
	}

}
