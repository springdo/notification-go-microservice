package main

import (
	"fmt"
	"net/http"

	"redhat/notification-microservice/api"
	"redhat/notification-microservice/config"
	"redhat/notification-microservice/domain"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	fmt.Println("Hello Notification Microservice :)")
	config := config.Read("config.toml")

	authEmailServer := domain.NewAuthEmailServer(config)

	r := mux.NewRouter()
	r.HandleFunc("/status", api.StatusHandler).Methods("GET")
	r.HandleFunc("/api/email", api.EmailHandler(authEmailServer)).Methods("POST")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
	})
	http.ListenAndServe(":8080", c.Handler(r))
}
