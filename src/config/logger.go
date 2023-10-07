package config

import (
	"log/slog"
	"os"
)

func initLogger() {
	/**
	* The slog package is a simple logging package that supports cross-platform color and concurrency.
	* slog is a simple logging package that supports cross-platform color and concurrency.
	**/
	opts := &slog.HandlerOptions{
		Level:     getLoggerLevel(),
		AddSource: true,
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	slog.SetDefault(logger)
}

func getLoggerLevel() slog.Level {
	logLevel := os.Getenv("LOG_LEVEL")
	switch logLevel {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
