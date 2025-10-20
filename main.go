package main

import (
	"context"
	"log"
	"log/slog"

	"github.com/renanmatosdacunha/golang-observability-otel.git/api"
	logs "github.com/renanmatosdacunha/golang-observability-otel.git/telemetry"
)

const (
	serverAddress = "0.0.0.0:8080"
)

func main() {
	baselogger := logs.NewLogger()
	logger := baselogger.With(
		slog.Group("Resource",
			slog.String("service.name", "golang-observability"),
			slog.String("deployment.environment", "delevopment"),
		),
	)

	server := api.NewServer(logger)
	logger.Log(context.Background(), slog.LevelInfo,
		"a iniciar o servidor", // Esta será a 'Body'
		// Adiciona os campos customizados para corresponder ao formato desejado.
		slog.Int("SeverityNumber", logs.GetSeverityNumber(slog.LevelInfo)),
		slog.Group("Attributes",
			slog.String("address", serverAddress),
		),
	)

	//logger.Info("Start Service", slog.String("address:", serverAddress))

	err := server.Start(serverAddress)
	if err != nil {
		log.Fatal("Cannot start server", err)
	}
}
