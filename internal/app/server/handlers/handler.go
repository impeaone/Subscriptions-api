package handlers

import (
	"agrigation_api/pkg/database/repository"
	"agrigation_api/pkg/logger/logger"
)

type Handler struct {
	db   *repository.Repository
	logs *logger.Log
}

func NewHandler(db *repository.Repository, logs *logger.Log) *Handler {
	return &Handler{db: db, logs: logs}
}
