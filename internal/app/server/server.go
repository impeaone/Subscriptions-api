package server

import (
	"agrigation_api/internal/app/server/handlers"
	"agrigation_api/internal/database/repository"
	"agrigation_api/internal/middleware"
	"agrigation_api/internal/service"
	"agrigation_api/pkg/config"
	logger2 "agrigation_api/pkg/logger"
	"context"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"sync"
)

type Server struct {
	Port        int
	Logger      logger2.MyLogger
	Router      http.Handler
	Postgres    *repository.Repository
	exitChan    chan struct{}
	connections *sync.WaitGroup
}

func NewServer(config *config.Config, logs logger2.MyLogger, service service.Subscriptions) *Server {
	port := config.Port

	router := http.NewServeMux()
	serverHandlers := handlers.NewHandler(service, logs)

	// Crud-операции
	router.HandleFunc("GET /api/v1/subscriptions/", serverHandlers.GetSubscription)
	router.HandleFunc("POST /api/v1/subscriptions/", serverHandlers.CreateSubscription)
	router.HandleFunc("PUT /api/v1/subscriptions/", serverHandlers.UpdateSubscription)
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

	exitChan := make(chan struct{})
	// Middleware
	shutdownMiddleware := middleware.ShutdownMiddleware(exitChan, router)
	loggerRouter := middleware.LoggerMiddleware(logs, shutdownMiddleware)
	PanicsRouter := middleware.PanicMiddleware(logs, loggerRouter)

	return &Server{
		Port:        port,
		Logger:      logs,
		Router:      PanicsRouter,
		exitChan:    exitChan,
		connections: &sync.WaitGroup{},
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	close(s.exitChan)

	finished := make(chan struct{})
	go func() {
		s.connections.Wait()
		close(finished)
	}()

	select {
	case <-finished:
		// Все операции завершилсь
		return nil
	case <-ctx.Done():
		// Время истекло
		return ctx.Err()
	}
}
