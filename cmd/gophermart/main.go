package main

import (
	"log"
	"net/http"
	"x2/cmd/handlers"

	"github.com/gorilla/mux"
)

func main() {
	st := handlers.Starts()
	r := mux.NewRouter()
	r.HandleFunc("/api/user/register", st.Register).Methods(http.MethodPost)
	r.HandleFunc("/api/user/login", st.Login).Methods(http.MethodPost)
	r.HandleFunc("/api/user/orders", st.UpOrder).Methods(http.MethodPost)
	r.HandleFunc("/api/user/orders", st.GetOrder).Methods(http.MethodGet)
	r.HandleFunc("/api/user/balance", st.Balance).Methods(http.MethodGet)
	r.HandleFunc("/api/user/balance/withdraw", st.LossBonus).Methods(http.MethodPost)
	r.HandleFunc("/api/user/withdrawals", st.Info).Methods(http.MethodGet)
	r.HandleFunc("/api/orders/{number}", st.InfoBonus).Methods(http.MethodGet)
	log.Println("server is running")
	err := http.ListenAndServe(st.URL, r)
	if err != nil {
		panic(err)
	}
	log.Println("server is running")
}
