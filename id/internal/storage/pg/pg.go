package pg

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/maximegorov13/chat-app/id/configs"
)

type Postgres struct {
	Db *sqlx.DB
}

func New(conf *configs.Config) (*Postgres, error) {
	db, err := sqlx.Connect("postgres", conf.Postgres.Url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return &Postgres{db}, nil
}
