package api

import (
	"errors"
	"log/slog"
	"net/http"
	"os"

	"github.com/FunnySneko/ReadPill/server/internal/aggregation"
	"github.com/FunnySneko/ReadPill/server/internal/db"
	"github.com/FunnySneko/ReadPill/server/internal/user_actions"
)

//go:generate go tool oapi-codegen -generate std-http,types -package api -o api.gn.go ../../../api/api.yml

var ErrUnauthorized = errors.New("unauthorized")

func NewServerHandler(database *db.Db, aggregator *aggregation.Aggregator, authHander *user_actions.ActionsHandler) *ServerHandler {
	serverHandler := ServerHandler{
		Db:     database,
		Agg:    aggregator,
		Ah:     authHander,
		secret: []byte(os.Getenv("JWTSECRET")),
	}
	return &serverHandler
}

type ServerHandler struct {
	Db     *db.Db
	Agg    *aggregation.Aggregator
	Ah     *user_actions.ActionsHandler
	secret []byte
}

func ErrorOut(w http.ResponseWriter, err error) {
	slog.Error(err.Error())

	var message string
	var statusCode int

	if errors.Is(err, ErrUnauthorized) {
		message = "unauthorized"
		statusCode = http.StatusUnauthorized
	} else if errors.Is(err, user_actions.ErrEmailTaken) {
		message = "email already taken"
		statusCode = http.StatusConflict
	} else if errors.Is(err, user_actions.ErrWrongUser) || errors.Is(err, user_actions.ErrWrongPassword) {
		message = "wrong email or password"
		statusCode = http.StatusUnauthorized
	} else if errors.Is(err, ErrInvalidRatings) {
		message = "invalid ratings"
		statusCode = http.StatusBadRequest
	} else {
		message = "internal server error"
		statusCode = http.StatusInternalServerError
	}

	http.Error(w, message, statusCode)
}
