package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/renanmatosdacunha/golang-observability-otel.git/api"
	logs "github.com/renanmatosdacunha/golang-observability-otel.git/telemetry"
)

const (
	serverAddress = "0.0.0.0:8080"
)

func main() {
	// Cria um contexto que lida com sinais de interrupção (Ctrl+C) para um shutdown gracioso.
	//ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	//defer cancel()

	// Inicializa o logger do OpenTelemetry através do nosso pacote de logs.
	//logger, shutdown := logs.InitOtelLogging(ctx)
	// Garante que a função de shutdown seja chamada no final para enviar todos os logs em buffer.
	//defer shutdown()
	////////////////////////////////////////SIGN/////////////////////////////////////////////////
	// Usa um contexto de fundo simples, uma vez que já não estamos a lidar com sinais.
	ctx := context.Background()
	logger, shutdownLogs := logs.InitOtelLogging(ctx)
	defer shutdownLogs()
	server := api.NewServer(logger)

	logMessage := fmt.Sprintf("Start Service: %s", serverAddress)
	logger.Info(logMessage)

	err := server.Start(serverAddress)
	if err != nil {
		logger.Error("não foi possível iniciar o servidor", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
