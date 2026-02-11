package handlers

import (
	"agrigation_api/internal/service"
	logger2 "agrigation_api/pkg/logger"
)

type Handler struct {
	serv service.Subscriptions
	logs logger2.MyLogger
}

func NewHandler(service service.Subscriptions, logs logger2.MyLogger) *Handler {
	return &Handler{serv: service, logs: logs}
}
