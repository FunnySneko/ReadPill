package user_actions

import "github.com/FunnySneko/ReadPill/server/internal/db"

func NewActionsHandler(database *db.Db) *ActionsHandler {
	return &ActionsHandler{
		db: database,
	}
}

type ActionsHandler struct {
	db *db.Db
}
