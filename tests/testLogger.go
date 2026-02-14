package tests

import (
	logger2 "agrigation_api/pkg/logger"
	"log"
)

type TestLog struct {
	Level string
}

func NewTestLog(level string) logger2.MyLogger {
	choseLevel := level
	if level != "WARNING" && level != "ERROR" {
		choseLevel = "INFO"
	}
	return &TestLog{
		Level: choseLevel,
	}
}

func (logs *TestLog) Info(message string, place string) {
	if logs.Level != "INFO" {
		return
	}
	log.Println("\nLevel: Info" + "\nMessage: " + message + "\nPlace: " + place + "\n")

}

func (logs *TestLog) Warning(message string, place string) {
	if logs.Level == "ERROR" {
		return
	}
	log.Println("\nLevel: Warning" + "\nMessage: " + message + "\nPlace: " + place + "\n")

}

func (logs *TestLog) Error(message string, place string) {
	log.Println("\nLevel: Error" + "\nMessage: " + message + "\nPlace: " + place + "\n")
}
