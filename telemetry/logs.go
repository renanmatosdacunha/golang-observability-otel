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

//func NewLogger() *slog.Logger {
//	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
//	return logger
//}

// NewLogger cria um logger slog altamente customizado para imitar o formato OTel.
func NewLogger() *slog.Logger {

	// O HandlerOptions permite-nos customizar o comportamento do handler de JSON.
	opts := slog.HandlerOptions{
		// A função ReplaceAttr é chamada para cada atributo antes de ser impresso.
		// É aqui que renomeamos as chaves e adicionamos novos campos.
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			// Renomeia 'time' para 'Timestamp' e usa o formato UnixNano.
			case slog.TimeKey:
				a.Key = "Timestamp"
				a.Value = slog.Int64Value(a.Value.Time().UnixNano())

			// Renomeia 'msg' para 'Body'.
			case slog.MessageKey:
				a.Key = "Body"

			// Renomeia 'level' para 'SeverityText'.
			case slog.LevelKey:
				a.Key = "SeverityText"
				// Também adicionamos o 'SeverityNumber' correspondente.
				// No entanto, ReplaceAttr só pode substituir um atributo por outro.
				// A melhor abordagem é adicionar o SeverityNumber no momento do log.
				// Veja o ficheiro user.go para um exemplo.
			}
			return a
		},
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &opts))
	return logger
}

// Mapeia os níveis do slog para os números de severidade do OpenTelemetry.
func GetSeverityNumber(level slog.Level) int {
	switch level {
	case slog.LevelDebug:
		return 5
	case slog.LevelInfo:
		return 9
	case slog.LevelWarn:
		return 13
	case slog.LevelError:
		return 17
	default:
		return 9 // Padrão para INFO
	}
}
