package handlers

import (
	"agrigation_api/pkg/database/migration"
	"agrigation_api/pkg/logger/logger"
)

type Handler struct {
	db   *migration.Repository
	logs *logger.Log
}

func NewHandler(db *migration.Repository, logs *logger.Log) *Handler {
	return &Handler{db: db, logs: logs}
}
