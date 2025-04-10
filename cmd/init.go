package main

import (
	"log"
	"log/slog"
	"os"
	"strings"
)

var cfg *config
var logger *slog.Logger
var collector *prometheusCollector

func init() {
	var err error
	var level slog.Level
	logLevel := strings.ToUpper(os.Getenv("LOG_LEVEL"))
	if logLevel == "" || logLevel == "WARNING" {
		logLevel = "WARN"
	}
	err = level.UnmarshalText([]byte(logLevel))
	if err != nil {
		panic(err)
	}
	logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))

	cfg, err = getConfig()
	if err != nil {
		panic(err)
	}

	if len(cfg.PrometheusPath) > 0 && cfg.PrometheusPath[0] != '/' {
		cfg.PrometheusPath = "/" + cfg.PrometheusPath
	}

	collector = newPrometheusCollector(cfg)
}

type serverErrorLogWriter struct{}

func (*serverErrorLogWriter) Write(p []byte) (int, error) {
	m := string(p)
	logger.Debug(m[:len(m)-1], "source", "http")
	return len(p), nil
}

func newServerErrorLog() *log.Logger {
	return log.New(&serverErrorLogWriter{}, "", 0)
}
