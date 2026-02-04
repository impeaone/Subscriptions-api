package server

import (
	"agrigation_api/internal/app/server/handlers"
	"agrigation_api/internal/middleware"
	"agrigation_api/pkg/config"
	"agrigation_api/pkg/database/postgres"
	"agrigation_api/pkg/database/redis"
	"agrigation_api/pkg/logger/logger"
	"net/http"
)

type Server struct {
	Port     int
	Logger   *logger.Log
	Router   http.Handler
	Postgres *postgres.Postgres
	Redis    *redis.Redis
}

func NewServer(config *config.Config, logs *logger.Log, pgs *postgres.Postgres) *Server {
	port := config.Port

	router := http.NewServeMux()

	// Crud-операции
	router.HandleFunc("GET /api/v1/subscription/{id}", handlers.GetSubscription)
	router.HandleFunc("POST /api/v1/subscription", handlers.CreateSubscription)
	router.HandleFunc("UPDATE /api/v1/subscription", handlers.UpdateSubscription)
	router.HandleFunc("DELETE /api/v1/subscription", handlers.DeleteSubscription)

	//health check
	router.HandleFunc("GET /health", handlers.HealthCheck)

	// Middleware
	loggerRouter := middleware.LoggerMiddleware(logs, router)
	PanicsRouter := middleware.PanicMiddleware(loggerRouter, logs)

	return &Server{
		Port:   port,
		Logger: logs,
		Router: PanicsRouter,
	}
}
