package logger

import (
	consts "agrigation_api/pkg/Constants"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

// TimeFormat - переменная для вывода времени в более понятном виде
var TimeFormat = time.Now().Format("02.01.2006 15:04:05")

// Log - структура для логера
type Log struct {
	Level string // Level - уровень логирования. INFO, WARNING, ERROR
}

// NewLog - конструктор для стуктуры лога
func NewLog(level string) *Log {
	choseLevel := level
	if level != "WARNING" && level != "ERROR" {
		choseLevel = "INFO"
	}
	return &Log{
		Level: choseLevel,
	}
}

// Info - метод для вывода логов с пометкой Info(обычные логи)
func (logs *Log) Info(message string, place string) {
	if logs.Level != "INFO" {
		return
	}
	log.Println("\nLevel: Info" + "\nMessage: " + message + "\nPlace: " + place + "\n")
	go WriteLogsToFile(TimeFormat +
		"\nLevel: Info" + "\nMessage: " + message + "\nPlace: " + place)
}

// Warning - метод для вывода логов с пометкой Warning(не крашат программу но опасны)
func (logs *Log) Warning(message string, place string) {
	if logs.Level == "ERROR" {
		return
	}
	log.Println("\nLevel: Warning" + "\nMessage: " + message + "\nPlace: " + place + "\n")
	go WriteLogsToFile(TimeFormat +
		"\nLevel: Warning" + "\nMessage: " + message + "\nPlace: " + place)
}

// Error - метод для вывода логов с пометкой Error (могут положить все)
func (logs *Log) Error(message string, place string) {
	log.Println("\nLevel: Error" + "\nMessage: " + message + "\nPlace: " + place + "\n")
	go WriteLogsToFile(TimeFormat +
		"\nLevel: Error" + "\nMessage: " + message + "\nPlace: " + place)
}

// GetPlace - функция для получения места вызова какой-то другой функции
func GetPlace() string {
	_, file, line, _ := runtime.Caller(1)
	split := strings.Split(file, "/")
	StartFile := split[len(split)-1]
	place := StartFile + ":" + strconv.Itoa(line)
	return place
}

// FileMTX - мютекс для записи в файл
var FileMTX sync.Mutex

// WriteLogsToFile - функция для записи лога в файл
func WriteLogsToFile(LogText string) {
	var logPath string
	if runtime.GOOS == "windows" {
		logPath = consts.LoggerPathWindows
	} else {
		logPath = consts.LoggerPathLinux
	}
	FileMTX.Lock()
	defer FileMTX.Unlock()
	file, err := os.OpenFile(
		logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("\nLevel: Error" + "\nMessage: " + consts.LogFileDoesNotOpen + ": " + err.Error() + "\nPlace: " +
			GetPlace() + "\n")
	}
	_, err = file.WriteString(LogText + "\n\n")
	if err != nil {
		log.Println("\nLevel: Error" + "\nMessage: " + consts.LogFileDoesNotWrite + ": " + err.Error() + "\nPlace: " +
			GetPlace() + "\n")
	}
	file.Close()
}
