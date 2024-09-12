package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"x2/internal/config"
	"x2/internal/handlers"
	"x2/internal/middleware"

	"github.com/gorilla/mux"
)

func main() {
	st := handlers.Init()
	r := mux.NewRouter()
	r.Use(middleware.AuthMiddleware)
	r.HandleFunc("/api/user/register", st.Register).Methods(http.MethodPost)
	r.HandleFunc("/api/user/login", st.Login).Methods(http.MethodPost)
	r.HandleFunc("/api/user/orders", st.UpOrder).Methods(http.MethodPost)
	r.HandleFunc("/api/user/orders", st.GetOrder).Methods(http.MethodGet)
	r.HandleFunc("/api/user/balance", st.Balance).Methods(http.MethodGet)
	r.HandleFunc("/api/user/balance/withdraw", st.LossBonus).Methods(http.MethodPost)
	r.HandleFunc("/api/user/withdrawals", st.Info).Methods(http.MethodGet)

	go func() {
		ticker := time.Tick(3 * time.Second)
		for range ticker {
			err := st.Bonus()
			if err != nil {
				//правильно логировать ошибки
				fmt.Println(err)
			}
		}
	}()

	log.Println("server is running")
	log.Fatal(http.ListenAndServe(config.GetServerAddress(), r))
	//панику лучше ну вызывать

}
