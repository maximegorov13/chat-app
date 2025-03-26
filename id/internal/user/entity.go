package user

import "time"

type User struct {
	ID        int64     `db:"id"`
	Login     string    `db:"login"`
	Name      string    `db:"name"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
