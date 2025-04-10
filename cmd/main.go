package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	broker, err := createBroker(cfg)
	if err != nil {
		logger.Error("failed to create MQTT broker", "error", err)
		return
	}
	broker.Subscribe("#", 1, shellyHandler)

	server := newServer(cfg)

	go func() {
		logger.Info("starting MQTT broker", "host", cfg.MqttHost, "port", cfg.MqttPort)
		if err := broker.Serve(); err != nil {
			logger.Error("failed to start MQTT broker", "error", err)
		}
	}()

	go func() {
		uri := fmt.Sprintf("%s://%s:%d%s", server.Scheme(), cfg.PrometheusHost, cfg.PrometheusPort, cfg.PrometheusPath)
		logger.Info("starting metrics server", "uri", uri)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			logger.Error("failed to start metrics server", "error", err)
		}
	}()

	<-ctx.Done()
	logger.Info("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("failed to shutdown metrics server", "error", err)
	}

	if err := broker.Close(); err != nil {
		logger.Error("failed to shutdown MQTT broker", "error", err)
	}

	logger.Info("shutdown complete")
}
