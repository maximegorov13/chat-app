package pg

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/maximegorov13/chat-app/id/configs"
)

type Postgres struct {
	Sqlx    *sqlx.DB
	Builder squirrel.StatementBuilderType
}

func New(conf *configs.Config) (*Postgres, error) {
	db, err := sqlx.Connect("postgres", conf.Postgres.Url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	
	return &Postgres{
		Sqlx:    db,
		Builder: builder,
	}, nil
}
