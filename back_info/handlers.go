package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
)

func EnergyTransactionHandler(conn *pgx.Conn, w http.ResponseWriter, r *http.Request) {

	et := &EnergyTransaction{}

	switch r.Method {
	case http.MethodGet:
		getEThandler(conn, w, r)
	case http.MethodPost:
		err := json.NewDecoder(r.Body).Decode(et)
		if err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		err1 := insert_energyTransaction(et, conn)
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)

	case http.MethodDelete:
		id := url_params(r.URL.Path, "/energy-transaction/")
		if id == "" {
			http.Error(w, "energy_transaction id is required", http.StatusBadRequest)
			return
		}
		err2 := delete_energyTransaction(id, conn)
		if err2 != nil {
			http.Error(w, err2.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	case http.MethodPut:
		err := json.NewDecoder(r.Body).Decode(et)
		if err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		err3 := update_energyTransaction(et, conn)
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

func getEThandler(conn *pgx.Conn, w http.ResponseWriter, r *http.Request) {

	url_params_p := url_params(r.URL.Path, "/energy-transaction/")
	if url_params_p == "" {
		energyTransactions, err := get_energyTransactions(conn)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(energyTransactions)

	} else {
		energyTransaction, err := get_energyTransactionsId(conn, url_params_p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(energyTransaction)

	}

}

func TakeOffHandler(conn *pgx.Conn, w http.ResponseWriter, r *http.Request) {

	to := &TakeOff{}

	switch r.Method {
	case http.MethodGet:
		getTOHandler(conn, w, r)
	case http.MethodPost:
		err := json.NewDecoder(r.Body).Decode(to)
		if err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		err1 := insert_takeoff(to, conn)
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)

	case http.MethodDelete:
		id := url_params(r.URL.Path, "/take-off/")
		if id == "" {
			http.Error(w, "take_off id is required", http.StatusBadRequest)
			return
		}
		err2 := delete_takeoff(id, conn)
		if err2 != nil {
			http.Error(w, err2.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	case http.MethodPut:
		err := json.NewDecoder(r.Body).Decode(to)
		if err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		err3 := update_takeoff(to, conn)
		if err3 != nil {
			http.Error(w, err3.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}

}

func getTOHandler(conn *pgx.Conn, w http.ResponseWriter, r *http.Request) {

	url_params_p := url_params(r.URL.Path, "/take-off/")
	if url_params_p == "" {
		take_offs, err := get_takeoffs(conn)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(take_offs)

	} else {
		take_off, err := get_takeoffId(conn, url_params_p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(take_off)

	}

}

func TenantHandler(conn *pgx.Conn, w http.ResponseWriter, r *http.Request) {

	tnt := &Tenant{}

	switch r.Method {
	case http.MethodGet:
		getTnHandler(conn, w, r)
	case http.MethodPost:
		err := json.NewDecoder(r.Body).Decode(tnt)
		if err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		err1 := insert_tenant(tnt, conn)
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)

	case http.MethodDelete:
		id := url_params(r.URL.Path, "/tenant/")
		if id == "" {
			http.Error(w, "tenant id is required", http.StatusBadRequest)
			return
		}
		err2 := delete_tenant(id, conn)
		if err2 != nil {
			http.Error(w, err2.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	case http.MethodPut:
		err := json.NewDecoder(r.Body).Decode(tnt)
		if err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		err3 := update_tenant(tnt, conn)
		if err3 != nil {
			http.Error(w, err3.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}

}

func getTnHandler(conn *pgx.Conn, w http.ResponseWriter, r *http.Request) {

	url_params_p := url_params(r.URL.Path, "/tenant/")
	if url_params_p == "" {
		take_offs, err := get_Tenants(conn)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(take_offs)

	} else {
		take_off, err := get_TenantId(conn, url_params_p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(take_off)

	}

}

func DestinationHandler(conn *pgx.Conn, w http.ResponseWriter, r *http.Request) {

	dtn := &Destination{}

	switch r.Method {
	case http.MethodGet:
		getDtnHandler(conn, w, r)
	case http.MethodPost:
		err := json.NewDecoder(r.Body).Decode(dtn)
		if err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		err1 := insert_destination(dtn, conn)
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)

	case http.MethodDelete:
		id := url_params(r.URL.Path, "/destination/")
		if id == "" {
			http.Error(w, "destination id is required", http.StatusBadRequest)
			return
		}
		err2 := delete_destination(id, conn)
		if err2 != nil {
			http.Error(w, err2.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	case http.MethodPut:
		err := json.NewDecoder(r.Body).Decode(dtn)
		if err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		err3 := update_destinantion(dtn, conn)
		if err3 != nil {
			http.Error(w, err3.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}

}

func getDtnHandler(conn *pgx.Conn, w http.ResponseWriter, r *http.Request) {

	url_params_p := url_params(r.URL.Path, "/destination/")
	if url_params_p == "" {
		take_offs, err := get_destinations(conn)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(take_offs)

	} else {
		destination, err := get_destinationsId(conn, url_params_p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(destination)

	}

}
