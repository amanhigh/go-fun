package telemetry

import (
	"bytes"
	"os"

	"github.com/amanhigh/go-fun/models/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const (
	logSeverityKey = attribute.Key("log.severity")
	logMessageKey  = attribute.Key("log.message")
)

// InitLogger initializes the logger with the specified level.
// It takes a parameter level of type zerolog.Level.
//
// level - zerolog.DebugLevel (Verbose) to zerolog.ErrorLevel (Limited), or zerolog.FatalLevel (Critical)
func InitLogger(config config.Log) {
	// Level
	zerolog.SetGlobalLevel(config.LogLevel)

	// Formatter
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006-01-02 15:04:05"}
	log.Logger = zerolog.New(output).With().Timestamp().Logger()

	// Add the OTEL Trace Hook
	log.Hook(NewZeroOtelHook())
}

func InitTestLogger(buffer *bytes.Buffer) {
	log.Logger = zerolog.New(buffer).With().Timestamp().Logger()
}

// ZeroOtelHook is a Zerolog hook that adds logs to the active span as events.
type ZeroOtelHook struct {
}

// NewZeroOtelHook returns a Zerolog hook.
func NewZeroOtelHook() zerolog.Hook {
	return &ZeroOtelHook{}
}

// Run adds trace context to the logger context.
func (h *ZeroOtelHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	ctx := e.GetCtx()
	if ctx == nil {
		return
	}

	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return
	}

	// Add events to the span
	attrs := make([]attribute.KeyValue, 0, 2)
	attrs = append(attrs, logSeverityKey.String(level.String()))
	attrs = append(attrs, logMessageKey.String(msg))

	span.AddEvent("log", trace.WithAttributes(attrs...))

	// Set status if level is error
	if level <= zerolog.ErrorLevel {
		span.SetStatus(codes.Error, msg)
	}
}
