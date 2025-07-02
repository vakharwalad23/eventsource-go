package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vakharwalad23/eventsource-starter-go/internal/app"
	"github.com/vakharwalad23/eventsource-starter-go/internal/domain"
)

type AccountService interface {
	CreateAccount(ctx context.Context, id string) error
	Deposit(ctx context.Context, id string, amount float64) error
	Withdraw(ctx context.Context, id string, amount float64) error
	GetAccount(ctx context.Context, id string) (*domain.Account, error)
}

func RgisterHandlers(r *mux.Router, svc *app.AccountService) {
	r.HandleFunc("/accounts", func(w http.ResponseWriter, r *http.Request) {
		var req struct{ ID string }
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := svc.CreateAccount(r.Context(), req.ID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}).Methods("POST")
	r.HandleFunc("/accounts/{id}/deposit", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		var req struct{ Amount float64 }
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := svc.Deposit(r.Context(), id, req.Amount); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}).Methods("POST")

	r.HandleFunc("/accounts/{id}/withdraw", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		var req struct{ Amount float64 }
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := svc.Withdraw(r.Context(), id, req.Amount); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}).Methods("POST")

	r.HandleFunc("/accounts/{id}/balance", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		acc, err := svc.GetAccount(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(acc)
	}).Methods("GET")

}
