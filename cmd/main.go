package main

import (
	"auth/pkg/config"
	"auth/service/handler"
	"auth/service/repo"
	"auth/service/usecase"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
	"log"
	"net/http"
)

func main() {
	connConfig, _ := pgx.ParseDSN(config.ConnInfo())
	conn, err := pgx.Connect(connConfig)
	if err != nil {
		log.Fatal("error:", err)
	}
	userR := repo.NewUserDB(conn)
	userS := usecase.NewUserService(userR)
	h := handler.NewHandler(userS)
	r := mux.NewRouter()
	r.HandleFunc("/auth", h.AuthHandle).Methods("POST")
	r.Handle("/hi", h.Authorization(http.HandlerFunc(h.Hello))).Methods("GET")
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
