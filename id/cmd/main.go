package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/maximegorov13/chat-app/id/configs"
	authhttp "github.com/maximegorov13/chat-app/id/internal/auth/delivery/http"
	authredis "github.com/maximegorov13/chat-app/id/internal/auth/repository/redis"
	authservice "github.com/maximegorov13/chat-app/id/internal/auth/service"
	"github.com/maximegorov13/chat-app/id/internal/storage/pg"
	"github.com/maximegorov13/chat-app/id/internal/storage/redis"
	userhttp "github.com/maximegorov13/chat-app/id/internal/user/delivery/http"
	userpg "github.com/maximegorov13/chat-app/id/internal/user/repository/pg"
	userservice "github.com/maximegorov13/chat-app/id/internal/user/service"
)

func main() {
	conf, err := configs.Load()
	if err != nil {
		log.Fatal(err)
	}

	pgClient, err := pg.New(conf)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := pgClient.Sqlx.DB.Close(); err != nil {
			log.Printf("Error when closing the Postgres connection: %v\n", err)
		}
	}()

	redisClient, err := redis.New(context.Background(), conf)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Printf("Error when closing the Redis connection: %v\n", err)
		}
	}()

	// Repositories
	userRepo := userpg.NewUserRepository(pgClient)
	tokenRepo := authredis.NewTokenRepository(redisClient)

	// Services
	userService := userservice.NewUserService(userservice.UserServiceDeps{
		UserRepo: userRepo,
	})
	authService := authservice.NewAuthService(authservice.AuthServiceDeps{
		UserRepo:  userRepo,
		TokenRepo: tokenRepo,
	})

	router := http.NewServeMux()

	// Handlers
	userhttp.NewUserHandler(router, userhttp.UserHandlerDeps{
		Conf:        conf,
		UserService: userService,
		TokenRepo:   tokenRepo,
	})
	authhttp.NewAuthHandler(router, authhttp.AuthHandlerDeps{
		Conf:        conf,
		AuthService: authService,
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%v", conf.Server.Port),
		Handler: router,
	}

	log.Printf("Starting server on port %v\n", conf.Server.Port)
	log.Fatal(server.ListenAndServe())
}
