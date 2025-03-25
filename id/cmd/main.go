package main

import (
	"fmt"
	"github.com/maximegorov13/chat-app/id/configs"
	"github.com/maximegorov13/chat-app/id/internal/storage/pg"
	"log"
	"net/http"
)

func main() {
	conf, err := configs.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	pgClient, err := pg.New(conf)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := pgClient.DB.Close(); err != nil {
			log.Printf("Error when closing the Postgres connection: %v\n", err)
		}
	}()

	router := http.NewServeMux()
	server := &http.Server{
		Addr:    fmt.Sprintf(":%v", conf.Server.Port),
		Handler: router,
	}
	log.Printf("Starting server on port %v\n", conf.Server.Port)
	log.Fatal(server.ListenAndServe())
}
