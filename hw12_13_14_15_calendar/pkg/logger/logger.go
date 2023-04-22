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
	case "INFO":
		logLevel = LogLevelInfo
	case "ERROR":
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

func (l *Logger) Infof(f string, args ...any) {
	if l.level == LogLevelInfo {
		fmt.Fprintf(os.Stdout, f, args...)
	}
}

func (l *Logger) Error(msg string) {
	fmt.Fprintln(os.Stderr, msg)
}

func (l *Logger) Errorf(f string, args ...any) {
	fmt.Fprintf(os.Stderr, f, args...)
}
