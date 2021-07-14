package main

import (
	"auth/pkg/config"
	"auth/pkg/connector_db"
	"auth/service/handler"
	"auth/service/repo"
	"auth/service/usecase"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	_ "github.com/lib/pq"
	"log"
	"os"
	"os/signal"
)

func main() {
	conn, err := connector_db.ConnPostger(config.ConnInfo())
	if err != nil {
		log.Fatal(err)
	}
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
	gh := handler.NewGinHandler(userS)
	middleware := handler.NewMiddleware(userS)

	gr := gin.Default()
	api := gr.Group("/v1")
	api.POST("/login", gh.Login)
	api.POST("/registration", gh.Registration)
	// api for user
	apiU := api.Group("")
	authU := apiU.Group("/", middleware.AuthorizationMiddleware())
	authU.GET("/hi", middleware.Authorize("resource", "read"), gh.Hello)
	authU.GET("/logout", middleware.Authorize("resource", "read"), gh.Logout)
	//api for admin
	apiA := api.Group("/admin")
	authA := apiA.Group("/", middleware.AuthorizationMiddleware())
	authA.GET("/hi", middleware.Authorize("resource1", "read"), gh.HelloA)
	authA.GET("/logout", middleware.Authorize("resource1", "read"), gh.Logout)

	go func() {
		if err := gr.Run(config.HostServer()); err != nil {
			log.Fatal(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
