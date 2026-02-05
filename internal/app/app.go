package app

import (
	"agrigation_api/internal/app/server"
	"agrigation_api/pkg/config"
	"agrigation_api/pkg/database/repository"
	"agrigation_api/pkg/logger/logger"
	"fmt"
	"net/http"
)

type App struct {
	fileServer *server.Server
}

func NewApp(config *config.Config, logger *logger.Log, pgs *repository.Repository) *App {
	fileServer := server.NewServer(config, logger, pgs)
	return &App{
		fileServer: fileServer,
	}
}

func (app *App) Start() error {
	app.fileServer.Logger.Info(fmt.Sprintf("Server listening on port %d", app.fileServer.Port), logger.GetPlace())
	err := http.ListenAndServe(fmt.Sprintf(":%d", app.fileServer.Port), app.fileServer.Router)
	return err
}
