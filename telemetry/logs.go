package logs

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

// InitOtelLogging configura o OpenTelemetry e retorna um logger slog e uma função de shutdown.
func InitOtelLogging(ctx context.Context) (*slog.Logger, func()) {
	// 1. Configura o Recurso (Resource)
	// O recurso descreve a aplicação que está gerando os logs.
	res, err := resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName("seu-servico-go"),
			semconv.ServiceVersion("v1.0.0"),
			semconv.DeploymentEnvironment("development"),
		),
	)
	if err != nil {
		panic(err)
	}

	// 2. Configura o Exportador
	// Por enquanto, vamos exportar para a saída padrão (console) no formato OTLP JSON.
	// Em produção, você trocaria por um exportador para o seu backend (Jaeger, Datadog, etc).
	logExporter, err := stdoutlog.New(
		stdoutlog.WithPrettyPrint(), // Deixa o JSON mais legível no console
	)
	if err != nil {
		panic(err)
	}

	// 3. Cria o LoggerProvider do OpenTelemetry
	// Este é o núcleo do SDK de logs do OTel.
	loggerProvider := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(logExporter)),
		sdklog.WithResource(res),
	)

	// 4. Define o LoggerProvider como global
	global.SetLoggerProvider(loggerProvider)

	// 5. Cria o nosso Handler customizado (a "ponte")
	handler := &OtelHandler{
		logger: loggerProvider.Logger("slog-handler"),
	}

	// 6. Retorna um logger slog que usa nosso handler e a função de shutdown
	logger := slog.New(handler)
	shutdown := func() {
		if err := loggerProvider.Shutdown(ctx); err != nil {
			logger.Error("falha ao desligar o logger provider do OTel", slog.String("error", err.Error()))
		}
	}

	return logger, shutdown
}

// OtelHandler é a nossa "ponte" que implementa a interface slog.Handler.
type OtelHandler struct {
	logger *sdklog.Logger
}

// Handle é chamado para cada log do slog e o converte para o formato OTel.
func (h *OtelHandler) Handle(ctx context.Context, rec slog.Record) error {
	logRecord := sdklog.Record{}
	logRecord.SetTimestamp(rec.Time)
	logRecord.SetBody(sdklog.StringValue(rec.Message))

	// Converte o nível do slog para a severidade do OTel
	switch rec.Level {
	case slog.LevelDebug:
		logRecord.SetSeverityNumber(sdklog.SeverityNumberDebug)
		logRecord.SetSeverityText("DEBUG")
	case slog.LevelInfo:
		logRecord.SetSeverityNumber(sdklog.SeverityNumberInfo)
		logRecord.SetSeverityText("INFO")
	case slog.LevelWarn:
		logRecord.SetSeverityNumber(sdklog.SeverityNumberWarn)
		logRecord.SetSeverityText("WARN")
	case slog.LevelError:
		logRecord.SetSeverityNumber(sdklog.SeverityNumberError)
		logRecord.SetSeverityText("ERROR")
	}

	// Adiciona os atributos do slog como atributos do OTel
	var attrs []attribute.KeyValue
	rec.Attrs(func(a slog.Attr) bool {
		attrs = append(attrs, attribute.String(a.Key, a.Value.String()))
		return true
	})
	logRecord.AddAttributes(attrs...)

	// Emite o log no formato OpenTelemetry
	h.logger.Emit(ctx, logRecord)
	return nil
}

// Outras funções necessárias para a interface slog.Handler
func (h *OtelHandler) Enabled(_ context.Context, _ slog.Level) bool { return true }

// WithAttrs cria um novo handler com atributos adicionados.
func (h *OtelHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	var newAttrs []attribute.KeyValue
	for _, a := range attrs {
		newAttrs = append(newAttrs, attribute.String(a.Key, a.Value.String()))
	}
	// Cria um novo logger com os atributos para manter a imutabilidade
	return &OtelHandler{
		logger: h.logger, // O logger em si pode ser reutilizado
	}
}
func (h *OtelHandler) WithGroup(name string) slog.Handler { return h }
