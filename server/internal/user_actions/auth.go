// █▀▀▀▀▀█ ▀ █ ▄ ▀ ▄ █▀▀▀▀▀█    PROJECT: ReadPill
// █ ███ █  ▀█▄▄█  ▀ █ ███ █    AUTHOR:  FunnySneko
// █ ▀▀▀ █  ▄ ▄  █▀▀ █ ▀▀▀ █
// ▀▀▀▀▀▀▀ █▄▀ █▄▀ ▀ ▀▀▀▀▀▀▀    © 2026
// ▀█▄█▀ ▀▄▄▀▄█▀▄█▄  ▀▄▄▄▄▄▀
// ▀▀▄█▀█▀▀ ▀ ██▀██  ▄█▄█▄ █
// █▄█▀ █▀▀▄▀█▀  ▄▄█▀▀▀▄▀▄▀
// █▀▀▄▄▀▀▄▀▀▀ ▀ ▀█▀█  ████▀
// ▀ ▀  ▀▀▀█▄ ▀▀█▀ █▀▀▀█▄█▄▀
// █▀▀▀▀▀█  ██▀▄▄ ▄█ ▀ █▄▀ ▀
// █ ███ █ █▀█▄▀▀▀████▀▀▀▄█▄
// █ ▀▀▀ █ ▄▀▀█ ▀▀▄ ▀█▀█▀███
// ▀▀▀▀▀▀▀ ▀▀▀▀▀   ▀▀   ▀  ▀

package user_actions

import (
	"context"
	"errors"

	"github.com/FunnySneko/ReadPill/server/internal/db"
	"golang.org/x/crypto/bcrypt"
)

var ErrEmailTaken = errors.New("email taken")
var ErrWrongUser = errors.New("wrong user")
var ErrWrongPassword = errors.New("wrong password")

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (ah *ActionsHandler) SignUp(ctx context.Context, username string, email string, password string) (int, error) {
	_, err := ah.db.FindUser(ctx, email)
	if err == nil {
		return 0, ErrEmailTaken
	} else if !errors.Is(err, db.ErrUserNotFound) {
		return 0, err
	}
	passwordHash, err := HashPassword(password)
	if err != nil {
		return 0, err
	}
	userId, err := ah.db.CreateUser(ctx, username, email, passwordHash)
	if err != nil {
		return 0, err
	}
	return userId, nil
}

func (ah *ActionsHandler) LogIn(ctx context.Context, email string, password string) (int, error) {
	u, err := ah.db.FindUser(ctx, email)
	if err != nil {
		if errors.Is(err, db.ErrUserNotFound) {
			return 0, ErrWrongUser
		}
		return 0, err
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) != nil {
		return 0, ErrWrongPassword
	} else {
		return u.ID, nil
	}
}
