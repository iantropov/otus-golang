package logger

import (
	"fmt"
	"os"
)

type LogLevel string

const (
	LogLevelInfo  = "info"
	LogLevelError = "error"
)

type Logger struct {
	level LogLevel
}

func New(level string) (*Logger, error) {
	var logLevel LogLevel
	switch level {
	case "info":
		logLevel = LogLevelInfo
	case "error":
		logLevel = LogLevelError
	default:
		return nil, fmt.Errorf("invalid log level: %s", level)
	}
	return &Logger{logLevel}, nil
}

func (l *Logger) Info(msg string) {
	if l.level == LogLevelInfo {
		fmt.Fprintln(os.Stdout, msg)
	}
}

func (l *Logger) Error(msg string) {
	fmt.Fprintln(os.Stderr, msg)
}
