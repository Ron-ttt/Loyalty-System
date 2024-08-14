package main

import (
	"net/http"

	"github.com/Ron-ttt/x2/cmd/handlers"

	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/api/user/register", handlers.Register()).Methods(http.MethodPost)
	r.HandleFunc("/api/user/login", handlers.Login()).Methods(http.MethodPost)
	r.HandleFunc("/api/user/orders", handlers.UpOrder()).Methods(http.MethodPost)
	r.HandleFunc("/api/user/orders", handlers.GetOrder()).Methods(http.MethodGet)
	r.HandleFunc("/api/user/balance", handlers.Balance()).Methods(http.MethodGet)
	r.HandleFunc("/api/user/balance/withdraw", handlers.LossBonus()).Methods(http.MethodPost)
	r.HandleFunc("/api/user/withdrawals", handlers.Info()).Methods(http.MethodGet)
	r.HandleFunc("/api/orders/{number}", handlers.InfoBonus()).Methods(http.MethodGet)

	err := http.ListenAndServe("http://localhost:8080", r)

	if err != nil {
		panic(err)
	}
}
