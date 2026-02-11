package main

import (
	_ "agrigation_api/docs"
	"agrigation_api/internal/app"
	"agrigation_api/internal/database/repository"
	"agrigation_api/migrations"
	"agrigation_api/pkg/config"
	"agrigation_api/pkg/logger/logger"
	"agrigation_api/pkg/tools"
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

/*
Требуемые переменные окружения (.env файл для Docker-compose):

	Runtime:
	tools.GetEnvAsInt("NUM_CPU", runtime.NumCPU())

	Сам сервер:
	tools.GetEnvAsInt("SERVER_PORT", 11682)
	tools.GetEnv("SERVER_IP", "127.0.0.1")

	Postgres
	tools.GetEnv("PG_USER", "postgres")
	tools.GetEnv("PG_PASSWORD", "postgres")
	tools.GetEnv("PG_HOST", "localhost")
	tools.GetEnvAsInt("PG_PORT", 5432)
	tools.GetEnv("PG_DATABASE", "aggregation")

	logger:
	tools.GetEnv("LOGGER", "INFO")
*/

// @title Subscription Management API
// @version 1.0
// @description API for managing user subscriptions with period-based calculations
// @BasePath /api/v1
// @schemes http
func main() {
	// Ограничение ресурсов
	runtime.GOMAXPROCS(tools.GetEnvAsInt("NUM_CPU", runtime.NumCPU()))

	// Logger
	logs := logger.NewLog(tools.GetEnv("LOGGER", "INFO"))

	// Migrate
	if errMigrate := migrations.CheckAndCreateTables(); errMigrate != nil {
		logs.Error("Error to init tables: "+errMigrate.Error(), logger.GetPlace())
		return
	}
	logs.Info("Init Database successful", logger.GetPlace())

	// Инициализация Postgres
	rep, errRep := repository.InitRepository()
	if errRep != nil {
		logs.Error(fmt.Sprintf("Ошибка инициализации PostgreSQL: %v", errRep), logger.GetPlace())
		return
	}
	logs.Info("Успешное подключение к PostgreSQL", logger.GetPlace())

	// Инициализация конфига
	conf, err := config.ReadConfig()
	if err != nil {
		logs.Error(fmt.Sprintf("Reading config file error: %v", err), logger.GetPlace())
		return
	}
	logs.Info("Успешная инициализация конфига", logger.GetPlace())

	// Инициализация сервера
	application := app.NewApp(conf, logs, rep)
	go func() {
		if errStart := application.Start(); errStart != nil {
			logs.Error(fmt.Sprintf("Server Start error: %v", errStart), logger.GetPlace())
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	sig := <-quit
	logs.Info(fmt.Sprintf("Received signal: %v", sig), logger.GetPlace())

	ctx, clos := context.WithTimeout(context.Background(), 30*time.Second)
	defer clos()
	if errShut := application.ShutDown(ctx); errShut != nil {
		logs.Error("Error graceful shutdown. Heavy stopping...", logger.GetPlace())
		os.Exit(1)
		return
	}
	rep.CloseConnection()
	logs.Info("Shutdown successful", logger.GetPlace())
}
