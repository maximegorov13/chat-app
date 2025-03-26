package main

import (
	"fmt"
	"github.com/maximegorov13/chat-app/id/configs"
	"github.com/maximegorov13/chat-app/id/internal/storage/pg"
	userhttp "github.com/maximegorov13/chat-app/id/internal/user/delivery/http"
	userpg "github.com/maximegorov13/chat-app/id/internal/user/repository/pg"
	userservice "github.com/maximegorov13/chat-app/id/internal/user/service"
	"log"
	"net/http"
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

	// Repositories
	userRepo := userpg.NewUserRepository(pgClient)

	// Services
	userService := userservice.NewUserService(userservice.UserServiceDeps{
		UserRepo: userRepo,
	})

	router := http.NewServeMux()

	// Handlers
	userhttp.NewUserHandler(router, userhttp.UserHandlerDeps{
		Conf:        conf,
		UserService: userService,
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%v", conf.Server.Port),
		Handler: router,
	}

	log.Printf("Starting server on port %v\n", conf.Server.Port)
	log.Fatal(server.ListenAndServe())
}
