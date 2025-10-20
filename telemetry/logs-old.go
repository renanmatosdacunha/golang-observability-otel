package logs

import (
	"log/slog"
	"os"
)

// NewLogger cria e retorna uma nova instância de slog.Logger configurada.
// Esta função centraliza a configuração do logger para toda a aplicação,
// garantindo que todos os logs tenham o mesmo formato.
// Configura o logger para escrever no formato JSON na saída padrão (console).
// O JSON é um formato estruturado, ideal para ser processado por
// sistemas de coleta e análise de logs.

func NewLogger() *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	return logger
}
