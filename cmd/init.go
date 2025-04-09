package main

import (
	"log/slog"
	"os"

	"github.com/caarlos0/env"
)

var cfg config
var logger *slog.Logger
var collector *prometheusCollector

func init() {
	var level slog.Level
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" || logLevel == "WARNING" {
		logLevel = "WARN"
	}
	if err := level.UnmarshalText([]byte(logLevel)); err != nil {
		panic(err)
	}
	logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))

	if err := env.Parse(&cfg); err != nil {
		logger.Error("failed to parse environment variables", "error", err)
		os.Exit(1)
	}

	collector = newPrometheusCollector(&cfg)
}
