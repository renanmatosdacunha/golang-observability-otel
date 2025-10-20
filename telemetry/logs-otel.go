package logs

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

// InitOtelLogging configura o OpenTelemetry para gerar logs no formato OTLP/JSON.
func InitOtelLogging(ctx context.Context) (*slog.Logger, func()) {
	// 1. Configura o Recurso (Resource) que descreve a sua aplicação.
	res, err := resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName("golang-observability"),
			semconv.DeploymentEnvironment("delevopment"),
		),
	)
	if err != nil {
		panic(err)
	}

	// 2. Configura o Exportador para a saída padrão.
	// Este exportador gera o formato OTLP/JSON exato que você deseja.
	logExporter, err := stdoutlog.New()
	if err != nil {
		panic(err)
	}

	// 3. Cria o LoggerProvider do OpenTelemetry.
	loggerProvider := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(logExporter)),
		sdklog.WithResource(res),
	)

	// 4. Define o LoggerProvider como global.
	global.SetLoggerProvider(loggerProvider)

	// 5. Cria a nossa "ponte" que converte logs do slog para o formato OTel.
	handler := &OtelHandler{
		logger: loggerProvider.Logger("slog-handler"),
	}

	// 6. Retorna um logger slog que usa a nossa ponte e a função de shutdown.
	logger := slog.New(handler)
	shutdown := func() {
		if err := loggerProvider.Shutdown(ctx); err != nil {
			logger.Error("falha ao desligar o logger provider do OTel", slog.String("error", err.Error()))
		}
	}

	return logger, shutdown
}

// OtelHandler implementa a interface slog.Handler.
type OtelHandler struct {
	logger log.Logger
}

// Handle converte cada log do slog para o formato OTel.
func (h *OtelHandler) Handle(ctx context.Context, rec slog.Record) error {
	record := sdklog.NewRecord(rec.Time, rec.Message)

	// Adiciona os atributos do slog como atributos do OTel.
	rec.Attrs(func(a slog.Attr) bool {
		record.AddAttributes(attribute.String(a.Key, a.Value.String()))
		return true
	})

	// Converte o nível do slog para a severidade do OTel.
	switch rec.Level {
	case slog.LevelDebug:
		record.SetSeverity(log.SeverityDebug)
		record.SetSeverityText("DEBUG")
	case slog.LevelInfo:
		record.SetSeverity(log.SeverityInfo)
		record.SetSeverityText("INFO")
	case slog.LevelWarn:
		record.SetSeverity(log.SeverityWarn)
		record.SetSeverityText("WARN")
	case slog.LevelError:
		record.SetSeverity(log.SeverityError)
		record.SetSeverityText("ERROR")
	}

	// Emite o log através do SDK do OpenTelemetry.
	h.logger.Emit(ctx, record)
	return nil
}

// Funções necessárias para a interface slog.Handler.
func (h *OtelHandler) Enabled(_ context.Context, _ slog.Level) bool { return true }
func (h *OtelHandler) WithAttrs(attrs []slog.Attr) slog.Handler     { return h }
func (h *OtelHandler) WithGroup(name string) slog.Handler           { return h }
