package configs

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	Auth     AuthConfig
}

type ServerConfig struct {
	Port string
}

type PostgresConfig struct {
	Url string
}

type RedisConfig struct {
	Url string
}

type AuthConfig struct {
	Secret string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	return &Config{
		Server: ServerConfig{
			Port: os.Getenv("PORT"),
		},
		Postgres: PostgresConfig{
			Url: os.Getenv("POSTGRES_URL"),
		},
		Redis: RedisConfig{
			Url: os.Getenv("REDIS_URL"),
		},
		Auth: AuthConfig{
			Secret: os.Getenv("SECRET"),
		},
	}, nil
}
