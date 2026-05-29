package api

import (
	"encoding/json"
	"errors"
	"io"
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
	if err != nil {
		return err
	}
	return nil
}

func (s *ServerHandler) PostSignup(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		ErrorOut(w, err, "internal server error", http.StatusInternalServerError)
		return
	}

	sc := SignupCredentials{}
	err = json.Unmarshal(body, &sc)
	if err != nil {
		ErrorOut(w, err, "internal server error", http.StatusInternalServerError)
		return
	}

	userId, err := s.Ah.SignUp(sc.Username, sc.Email, sc.Password)
	if err != nil {
		if errors.Is(err, user_actions.ErrEmailTaken) {
			ErrorOut(w, err, "email already taken", http.StatusConflict)
		} else {
			ErrorOut(w, err, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	err = s.SetJWTToken(w, userId)
	if err != nil {
		ErrorOut(w, err, "internal server error", http.StatusInternalServerError)
		return
	}

	err = WriteJSON(w, http.StatusOK, map[string]string{
		"message": "signup successful",
	})
	if err != nil {
		ErrorOut(w, err, "internal server error", http.StatusInternalServerError)
	}
}

func (s *ServerHandler) PostLogin(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		ErrorOut(w, err, "internal server error", http.StatusInternalServerError)
		return
	}

	lc := LoginCredentials{}
	err = json.Unmarshal(body, &lc)
	if err != nil {
		ErrorOut(w, err, "internal server error", http.StatusInternalServerError)
		return
	}

	userId, err := s.Ah.LogIn(lc.Email, lc.Password)
	if err != nil {
		if errors.Is(err, user_actions.ErrWrongUser) || errors.Is(err, user_actions.ErrWrongPassword) {
			ErrorOut(w, err, "wrong email or password", http.StatusUnauthorized)
		} else {
			ErrorOut(w, err, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	err = s.SetJWTToken(w, userId)
	if err != nil {
		ErrorOut(w, err, "internal server error", http.StatusInternalServerError)
		return
	}

	err = WriteJSON(w, http.StatusOK, map[string]string{
		"message": "login successful",
	})
	if err != nil {
		ErrorOut(w, err, "internal server error", http.StatusInternalServerError)
	}
}
