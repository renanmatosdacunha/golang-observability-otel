package logs

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0" // ATUALIZADO para uma versão mais recente
)

// InitOtelLogging configura o OpenTelemetry para gerar logs no formato OTLP/JSON.
func InitOtelLogging(ctx context.Context) (*slog.Logger, func()) {
	// --- CORREÇÃO APLICADA AQUI ---
	// Em vez de usar resource.Merge, criamos um novo recurso do zero.
	// Isto também deteta os atributos padrão, mas de uma forma que evita conflitos de esquema.
	res, err := resource.New(ctx,
		resource.WithSchemaURL(semconv.SchemaURL),
		resource.WithAttributes(
			semconv.ServiceName("golang-observability"),
			semconv.DeploymentEnvironmentNameKey.String("development"), // Corrigido de "delevopment"
		),
	)
	if err != nil {
		panic(err)
	}

	// Configura o Exportador para a saída padrão.
	logExporter, err := stdoutlog.New()
	if err != nil {
		panic(err)
	}

	// Cria o LoggerProvider do OpenTelemetry.
	loggerProvider := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(logExporter)),
		sdklog.WithResource(res),
	)

	// Define o LoggerProvider como global.
	global.SetLoggerProvider(loggerProvider)

	// Cria a nossa "ponte" que converte logs do slog para o formato OTel.
	handler := &OtelHandler{
		logger: loggerProvider.Logger("slog-handler"),
	}

	// Retorna um logger slog que usa a nossa ponte e a função de shutdown.
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
	var record log.Record
	record.SetTimestamp(rec.Time)
	record.SetBody(log.StringValue(rec.Message))

	// Adiciona os atributos do slog como atributos do OTel.
	rec.Attrs(func(a slog.Attr) bool {
		record.AddAttributes(log.String(a.Key, a.Value.String()))
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
