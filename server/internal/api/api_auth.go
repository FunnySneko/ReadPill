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

package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/FunnySneko/ReadPill/server/internal/user_actions"
	"github.com/golang-jwt/jwt/v5"
)

func (s *ServerHandler) SetJWTToken(w http.ResponseWriter, userId int) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(30 * 24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(s.secret)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
	})
	return nil
}

func WriteJSON(w http.ResponseWriter, status int, payload any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(payload)
	return err
}

func (s *ServerHandler) PostSignup(w http.ResponseWriter, r *http.Request) {
	sc := SignupCredentials{}
	err := json.NewDecoder(r.Body).Decode(&sc)
	if err != nil {
		ErrorOut(w, err)
		return
	}

	userId, err := s.Ah.SignUp(r.Context(), sc.Username, sc.Email, sc.Password)
	if err != nil {
		if errors.Is(err, user_actions.ErrEmailTaken) {
			ErrorOut(w, err)
		} else {
			ErrorOut(w, err)
		}
		return
	}

	err = s.SetJWTToken(w, userId)
	if err != nil {
		ErrorOut(w, err)
		return
	}

	err = WriteJSON(w, http.StatusOK, map[string]string{
		"message": "signup successful",
	})
	if err != nil {
		ErrorOut(w, err)
	}
}

func (s *ServerHandler) PostLogin(w http.ResponseWriter, r *http.Request) {
	lc := LoginCredentials{}
	err := json.NewDecoder(r.Body).Decode(&lc)
	if err != nil {
		ErrorOut(w, err)
		return
	}

	userId, err := s.Ah.LogIn(r.Context(), lc.Email, lc.Password)
	if err != nil {
		if errors.Is(err, user_actions.ErrWrongUser) || errors.Is(err, user_actions.ErrWrongPassword) {
			ErrorOut(w, err)
		} else {
			ErrorOut(w, err)
		}
		return
	}

	err = s.SetJWTToken(w, userId)
	if err != nil {
		ErrorOut(w, err)
		return
	}

	err = WriteJSON(w, http.StatusOK, map[string]string{
		"message": "login successful",
	})
	if err != nil {
		ErrorOut(w, err)
	}
}
