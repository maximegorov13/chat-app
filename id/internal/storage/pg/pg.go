package pg

import (
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/maximegorov13/chat-app/id/configs"
)

type Postgres struct {
	Sqlx *sqlx.DB
	Sb   squirrel.StatementBuilderType
}

func NewPostgres(conf *configs.Config) (*Postgres, error) {
	db, err := sqlx.Connect("postgres", conf.Postgres.Url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	return &Postgres{
		Sqlx: db,
		Sb:   builder,
	}, nil
}
