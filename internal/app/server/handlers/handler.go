package handlers

import (
	"agrigation_api/internal/service"
	"agrigation_api/pkg/logger/logger"
)

type Handler struct {
	serv service.Subscriptions
	logs *logger.Log
}

func NewHandler(service service.Subscriptions, logs *logger.Log) *Handler {
	return &Handler{serv: service, logs: logs}
}
