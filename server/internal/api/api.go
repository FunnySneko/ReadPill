package api

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/FunnySneko/ReadPill/server/internal/aggregation"
	"github.com/FunnySneko/ReadPill/server/internal/db"
	"github.com/FunnySneko/ReadPill/server/internal/user_actions"
)

//go:generate go tool oapi-codegen -generate std-http,types -package api -o api.gn.go ../../../api/api.yml

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

func ErrorOut(w http.ResponseWriter, err error, message string, statusCode int) {
	slog.Error(err.Error())
	http.Error(w, message, statusCode)
}
