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

	mux.HandleFunc("/heartbeat", withInstanceHeader(heartbeatHandler))

	err2 := http.ListenAndServe(":8081", mux)
	if err2 != nil {
		log.Fatal(err2)
	}

}
