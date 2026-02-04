package main

import (
	"agrigation_api/internal/app"
	"agrigation_api/pkg/config"
	"agrigation_api/pkg/database/postgres"
	"agrigation_api/pkg/logger/logger"
	"agrigation_api/pkg/tools"
	"fmt"
	"runtime"
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
	tools.GetEnv("PG_DATABASE", "agrigations")

	logger:
	tools.GetEnv("CloudStorage_LOGGER", "INFO")

	REDIS:
	tools.GetEnv("REDIS_HOST", "localhost")
	tools.GetEnvAsInt("REDIS_PORT", 6379)
	tools.GetEnv("REDIS_PASSWORD", "")
	rdsDB := tools.GetEnvAsInt("REDIS_DB", 0)
*/
func main() {
	// Ограничение ресурсов
	runtime.GOMAXPROCS(tools.GetEnvAsInt("NUM_CPU", runtime.NumCPU()))

	// Logger
	logs := logger.NewLog(tools.GetEnv("CloudStorage_LOGGER", "INFO"))

	// Инициализация Postgres
	pgs, errPGS := postgres.InitPostgres()
	if errPGS != nil {
		logs.Error(fmt.Sprintf("Ошибка инициализации PostgreSQL: %v", errPGS), logger.GetPlace())
		return
	}
	logs.Info("Успешная инициализация PostgreSQL", logger.GetPlace())

	/*
		rds, errRds := redis.NewRedis()
		if errRds != nil {
			logs.Error(fmt.Sprintf("Ошибка инициализации Redis: %v", errRds), logger.GetPlace())
			return
		}
		logs.Info("Успешная инициализация Redis", logger.GetPlace())?
	*/

	// Инициализация конфига
	conf, err := config.ReadConfig()
	if err != nil {
		logs.Error(fmt.Sprintf("Reading config file error: %v", err), logger.GetPlace())
		return
	}
	logs.Info("Успешная инициализация конфига", logger.GetPlace())

	// Инициализация сервера
	if errStart := app.NewApp(conf, logs, pgs).Start(); errStart != nil {
		logs.Error(fmt.Sprintf("Server Start error: %v", errStart), logger.GetPlace())
		return
	}
}
