package main

import (
	"auth/pkg/config"
	"auth/pkg/connector_db"
	"auth/service/handler"
	"auth/service/repo"
	"auth/service/usecase"
	"context"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	conn := connector_db.ConnPostger(config.ConnInfo())
	connRedis := connector_db.ConnRedis(&redis.Options{
		Addr:     config.RedisAddr(),
		Password: config.RedisPass(),
		DB:       config.RedisDB(),
	})
	defer func() {
		conn.Close()
		_ = connRedis.Close()
	}()
	userR := repo.NewUserDB(conn, connRedis)
	userKR := repo.NewUserKR(connRedis)
	jwtU := repo.NewJwtRepo(connRedis)
	userS := usecase.NewUserService(userR, userKR, jwtU)
	h := handler.NewHandler(userS)
	r := mux.NewRouter()
	r.HandleFunc("/registration", h.RegHandler).Methods("POST")
	r.HandleFunc("/auth", h.AuthHandle).Methods("POST")
	r.Handle("/hi", h.AuthorizationMidlaware(http.HandlerFunc(h.Hello))).Methods("GET")
	http.Handle("/", r)
	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8080",
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Print("HTTP server error:", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Попытка корректного завершения
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}
