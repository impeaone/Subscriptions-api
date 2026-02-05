package server

import (
	"agrigation_api/internal/app/server/handlers"
	"agrigation_api/internal/middleware"
	"agrigation_api/pkg/config"
	"agrigation_api/pkg/database/repository"
	"agrigation_api/pkg/logger/logger"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

type Server struct {
	Port     int
	Logger   *logger.Log
	Router   http.Handler
	Postgres *repository.Repository
}

func NewServer(config *config.Config, logs *logger.Log, pgs *repository.Repository) *Server {
	port := config.Port

	router := http.NewServeMux()
	serverHandlers := handlers.NewHandler(pgs, logs)

	// Crud-операции
	router.HandleFunc("GET /api/v1/subscriptions/", serverHandlers.GetSubscription)
	router.HandleFunc("POST /api/v1/subscriptions/", serverHandlers.UpsertSubscription)
	router.HandleFunc("GET /api/v1/subscriptions/user/{id}", serverHandlers.ListUserSubscriptions)
	router.HandleFunc("DELETE /api/v1/subscriptions/", serverHandlers.DeleteSubscription)

	router.HandleFunc("GET /api/v1/subscriptions/total/", serverHandlers.CalculateTotalHandler)

	// health check
	router.HandleFunc("GET /health", serverHandlers.HealthCheck)

	// Swagger
	router.Handle("GET /swagger/", httpSwagger.WrapHandler)
	// Редирект с корня на Swagger UI
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/index.html", http.StatusFound)
	})

	// Middleware
	loggerRouter := middleware.LoggerMiddleware(logs, router)
	PanicsRouter := middleware.PanicMiddleware(loggerRouter, logs)

	return &Server{
		Port:   port,
		Logger: logs,
		Router: PanicsRouter,
	}
}
