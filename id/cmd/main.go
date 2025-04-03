package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/maximegorov13/chat-app/id/pkg/jwt"

	"github.com/maximegorov13/chat-app/id/configs"
	authhttp "github.com/maximegorov13/chat-app/id/internal/auth/delivery/http"
	authredis "github.com/maximegorov13/chat-app/id/internal/auth/repository/redis"
	authservice "github.com/maximegorov13/chat-app/id/internal/auth/service"
	"github.com/maximegorov13/chat-app/id/internal/keyreader"
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

	pgClient, err := pg.NewPostgres(conf)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := pgClient.Sqlx.Close(); err != nil {
			log.Printf("Error when closing the Postgres connection: %v\n", err)
		}
	}()

	redisClient, err := redis.NewRedis(context.Background(), conf)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Printf("Error when closing the Redis connection: %v\n", err)
		}
	}()

	keyReader := keyreader.NewKeyReader(conf.Auth.SecretKeysPath)
	privateKey, err := keyReader.ReadPrivateKey(conf.Auth.PostfixKeyAuth)
	if err != nil {
		log.Fatal(err)
	}
	publicKey, err := keyReader.ReadPublicKey(conf.Auth.PostfixKeyAuth)
	if err != nil {
		log.Fatal(err)
	}

	jwtMaker := jwt.NewJWT(privateKey, publicKey)

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
		JWT:       jwtMaker,
	})

	router := http.NewServeMux()

	// Handlers
	userhttp.NewUserHandler(router, userhttp.UserHandlerDeps{
		Conf:        conf,
		UserService: userService,
		TokenRepo:   tokenRepo,
		JWT:         jwtMaker,
	})
	authhttp.NewAuthHandler(router, authhttp.AuthHandlerDeps{
		Conf:        conf,
		AuthService: authService,
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%v", conf.Server.Port),
		Handler: router,
	}

	go func() {
		log.Printf("Starting server on port %v", conf.Server.Port)
		if err = server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
		log.Println("Stopped serving new connections")
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}
	log.Println("Graceful shutdown complete")
}
