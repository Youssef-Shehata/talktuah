package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

type logger struct {
	file *os.File
}
type LogLevel string

// Define constants for valid log levels
const (
	ERROR LogLevel = "ERROR"
	INFO  LogLevel = "INFO"
)

func ValidateLogLevel(level LogLevel) bool {
	switch level {
	case ERROR, INFO:
		return true
	default:
		return false
	}
}

func NewLogger(filePath string) (*logger, error) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	return &logger{
		file: file,
	}, nil
}

func (l *logger) Log(lvl LogLevel , err error ) {
	if !ValidateLogLevel(lvl) {
		if os.Getenv("PLATFORM") == "dev" {
			log.Fatalf("Invalid log level: %s", lvl)
		}
	}
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf(" [%s] %v : %v\n", timestamp, lvl, err)

	fmt.Println(msg)

	if _, writeErr := l.file.WriteString(msg); writeErr != nil {
		log.Printf("Failed to write to log file: %v\n", writeErr)
	}
}

func (l *logger) Close() error {
	return l.file.Close()
}
