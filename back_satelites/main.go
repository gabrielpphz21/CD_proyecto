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

	mux.HandleFunc("/heartbeat", withInstanceHeader(heartbeatHandler))

	err1 := http.ListenAndServe(":8080", mux)
	if err1 != nil {
		log.Fatal(err1)
	}

}
