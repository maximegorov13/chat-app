package configs

import (
	"fmt"
	"os"

	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	Auth     AuthConfig
}

func (c Config) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Server),
		validation.Field(&c.Postgres),
		validation.Field(&c.Redis),
		validation.Field(&c.Auth),
	)
}

type ServerConfig struct {
	Port string
}

func (s ServerConfig) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Port, validation.Required, is.Port),
	)
}

type PostgresConfig struct {
	Url string
}

func (p PostgresConfig) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Url, validation.Required, is.URL),
	)
}

type RedisConfig struct {
	Url string
}

func (r RedisConfig) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Url, validation.Required, is.URL),
	)
}

type AuthConfig struct {
	SecretKeysPath string
	PostfixKeyAuth string
}

func (a AuthConfig) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.SecretKeysPath, validation.Required),
		validation.Field(&a.PostfixKeyAuth, validation.Required),
	)
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	conf := &Config{
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
			SecretKeysPath: os.Getenv("SECRET_KEYS_PATH"),
			PostfixKeyAuth: os.Getenv("POSTFIX_KEY_AUTH"),
		},
	}

	if err := conf.Validate(); err != nil {
		return nil, fmt.Errorf("error validating config: %w", err)
	}

	return conf, nil
}
