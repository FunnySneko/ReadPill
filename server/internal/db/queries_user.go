package db

import (
	"context"
	"errors"

	"github.com/FunnySneko/ReadPill/server/internal/user"
	"github.com/jackc/pgx/v5"
)

var ErrUserNotFound = errors.New("user not found")

func (d *Db) FindUser(email string) (user.User, error) {
	u := user.User{}
	err := d.conn.QueryRow(context.Background(), `SELECT id, username, email, password_hash FROM "user" WHERE email = $1`, email).Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash)
	if errors.Is(err, pgx.ErrNoRows) {
		return u, ErrUserNotFound
	}
	return u, err
}

func (d *Db) CreateUser(username string, email string, passwordHash string) (int, error) {
	userId := 0
	err := d.conn.QueryRow(context.Background(), `INSERT INTO "user" (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id`, username, email, passwordHash).Scan(&userId)
	return userId, err
}
