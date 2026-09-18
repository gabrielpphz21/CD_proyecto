package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
)

type SateliteInfo struct {
	Satelite   Satelite        `json:"satelite"`
	H_satelite []HyperSatelite `json:"h_satelite"`
}

func SateliteHandler(conn *pgx.Conn, w http.ResponseWriter, r *http.Request) {

	stl := &Satelite{}

	switch r.Method {
	case http.MethodGet:
		rockets, err := get_satelites(conn)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(rockets)
	case http.MethodPost:
		err := json.NewDecoder(r.Body).Decode(stl)
		if err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		err1 := insert_satelite(stl, conn)
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)

	case http.MethodDelete:
		id := url_params(r.URL.Path, "/satelite/")
		if id == "" {
			http.Error(w, "rocket id is required", http.StatusBadRequest)
			return
		}
		err2 := delete_satelite(id, conn)
		if err2 != nil {
			http.Error(w, err2.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	case http.MethodPut:
		err := json.NewDecoder(r.Body).Decode(stl)
		if err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		err3 := update_satelite(stl, conn)
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

func HyperSateliteHandler(conn *pgx.Conn, w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		id := url_params(r.URL.Path, "/hyper-satelite/")
		if id == "" {
			http.Error(w, "no id provided", http.StatusBadRequest)
			return
		}
		si := &SateliteInfo{}
		stl, hs, err := satelite_info(id, conn)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		si.Satelite = *stl
		si.H_satelite = hs

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(si)

	case http.MethodPost:
		hr := &HyperSatelite{}
		err := json.NewDecoder(r.Body).Decode(hr)
		if err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		err1 := insert_hypersatelites(hr, conn)
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)

	case http.MethodDelete:
		id := url_params(r.URL.Path, "/hyper-satelite/")
		if id == "" {
			http.Error(w, "h_satelite id is required", http.StatusBadRequest)
			return
		}
		err2 := delete_hypersateliteLatest(id, conn)
		if err2 != nil {
			http.Error(w, err2.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}

}
