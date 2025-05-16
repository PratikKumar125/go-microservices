package logging

import (
	"os"
	"log/slog"
)

type Logger struct {
	slog *slog.Logger
}

func NewLogger(service, level string) *Logger {
	var lvl slog.Level
	switch level {
	case "debug":
		lvl = slog.LevelDebug
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl})
	return &Logger{slog: slog.New(handler).With("service", service)}
}

func (l *Logger) Info(msg string, args ...interface{}) {
	l.slog.Info(msg, args...)
}

func (l *Logger) Error(msg string, args ...interface{}) {
	l.slog.Error(msg, args...)
}