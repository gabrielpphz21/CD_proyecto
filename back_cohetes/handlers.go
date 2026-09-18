package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
)

type RocketInfo struct {
	Rocket    Rocket        `json:"rocket"`
	H_rokcets []HyperRocket `json:"h_rockets"`
}

func rocket_Handler(conn *pgx.Conn, w http.ResponseWriter, r *http.Request) {

	rkt := &Rocket{}

	switch r.Method {
	case http.MethodGet:
		rockets, err := get_rockets(conn)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(rockets)
	case http.MethodPost:
		err := json.NewDecoder(r.Body).Decode(rkt)
		if err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		err1 := insert_rocket(rkt, conn)
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)

	case http.MethodDelete:
		id := url_params(r.URL.Path, "/rocket/")
		if id == "" {
			http.Error(w, "rocket id is required", http.StatusBadRequest)
			return
		}
		err2 := delete_rocket(id, conn)
		if err2 != nil {
			http.Error(w, err2.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	case http.MethodPut:
		err := json.NewDecoder(r.Body).Decode(rkt)
		if err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		err3 := update_rocket(rkt, conn)
		if err3 != nil {
			http.Error(w, err3.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}

}

func url_params(url string, prefix string) string {
	_, after, found := strings.Cut(url, prefix)
	if found != false {
		return after
	}
	return ""
}

func HyperRocketHandler(conn *pgx.Conn, w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		id := url_params(r.URL.Path, "/hyper-rocket/")
		if id == "" {
			http.Error(w, "no id provided", http.StatusBadRequest)
			return
		}
		si := &RocketInfo{}
		rkt, hr, err := rocket_info(id, conn)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		si.Rocket = *rkt
		si.H_rokcets = hr

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(si)

	case http.MethodPost:
		hr := &HyperRocket{}
		err := json.NewDecoder(r.Body).Decode(hr)
		if err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		err1 := insert_hyperrocket(hr, conn)
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)

	case http.MethodDelete:
		id := url_params(r.URL.Path, "/hyper-rocket/")
		if id == "" {
			http.Error(w, "h_rocket id is required", http.StatusBadRequest)
			return
		}
		err2 := delete_hyperrocketLatest(id, conn)
		if err2 != nil {
			http.Error(w, err2.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}

}
