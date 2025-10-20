package main

import (
	"log"
	"log/slog"

	"github.com/renanmatosdacunha/golang-observability-otel.git/api"
	logs "github.com/renanmatosdacunha/golang-observability-otel.git/telemetry"
)

const (
	serverAddress = "0.0.0.0:8080"
)

func main() {
	logger := logs.NewLogger()
	server := api.NewServer(logger)
	logger.Info("Start Service", slog.String("address:", serverAddress))

	err := server.Start(serverAddress)
	if err != nil {
		log.Fatal("Cannot start server", err)
	}
}
