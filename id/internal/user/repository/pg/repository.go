package pg

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Masterminds/squirrel"

	"github.com/maximegorov13/chat-app/id/internal/storage/pg"
	"github.com/maximegorov13/chat-app/id/internal/user"
)

type UserRepository struct {
	db *pg.Postgres
}

func NewUserRepository(db *pg.Postgres) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (repo *UserRepository) Create(ctx context.Context, user *user.User) error {
	query, args, err := repo.db.Sb.
		Insert("users").
		Columns("login", "name", "password").
		Values(user.Login, user.Name, user.Password).
		Suffix("RETURNING id, created_at, updated_at").
		ToSql()
	if err != nil {
		return err
	}

	return repo.db.Sqlx.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (repo *UserRepository) FindByLogin(ctx context.Context, login string) (*user.User, error) {
	query, args, err := repo.db.Sb.
		Select("*").
		From("users").
		Where(squirrel.Eq{
			"login": login,
		}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var u user.User
	if err = repo.db.Sqlx.GetContext(ctx, &u, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &u, nil
}

func (repo *UserRepository) FindByID(ctx context.Context, id int64) (*user.User, error) {
	query, args, err := repo.db.Sb.
		Select("*").
		From("users").
		Where(squirrel.Eq{
			"id": id,
		}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var u user.User
	if err = repo.db.Sqlx.GetContext(ctx, &u, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &u, nil
}

func (repo *UserRepository) Update(ctx context.Context, user *user.User) error {
	query, args, err := repo.db.Sb.
		Update("users").
		SetMap(map[string]any{
			"login":      user.Login,
			"name":       user.Name,
			"password":   user.Password,
			"updated_at": squirrel.Expr("CURRENT_TIMESTAMP"),
		}).
		Where(squirrel.Eq{
			"id": user.ID,
		}).
		Suffix("RETURNING id, created_at, updated_at").
		ToSql()
	if err != nil {
		return err
	}

	return repo.db.Sqlx.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}
