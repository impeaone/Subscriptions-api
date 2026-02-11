package logger

import "agrigation_api/pkg/logger/logger"

type MyLogger interface {
	Info(string, string)
	Warning(string, string)
	Error(string, string)
}

func NewMyLogger(level string) MyLogger {
	return logger.NewLog(level)
}
