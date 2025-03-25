package configs

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	Server   ServerConfig
	Postgres PostgresConfig
}

type ServerConfig struct {
	Port string
}

type PostgresConfig struct {
	Url string
}

func LoadConfig() (*Config, error) {
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
	}, nil
}
