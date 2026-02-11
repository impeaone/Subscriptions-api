package app

import (
	"agrigation_api/internal/app/server"
	"agrigation_api/internal/service"
	"agrigation_api/pkg/config"
	"agrigation_api/pkg/logger/logger"
	"context"
	"fmt"
	"net/http"
)

type App struct {
	fileServer *server.Server
}

func NewApp(config *config.Config, logger *logger.Log, service service.Subscriptions) *App {
	fileServer := server.NewServer(config, logger, service)
	return &App{
		fileServer: fileServer,
	}
}

func (app *App) Start() error {
	app.fileServer.Logger.Info(fmt.Sprintf("Server listening on port %d", app.fileServer.Port), logger.GetPlace())
	err := http.ListenAndServe(fmt.Sprintf(":%d", app.fileServer.Port), app.fileServer.Router)
	return err
}

func (app *App) ShutDown(ctx context.Context) error {
	//defer app.fileServer.Postgres.CloseConnection()
	if err := app.fileServer.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
