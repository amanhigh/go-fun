package metrics

import (
	"context"

	"go.opentelemetry.io/otel"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func InitStdoutTracerProvider() {
	exporter, _ := stdout.New(stdout.WithPrettyPrint())
	traceprovider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSyncer(exporter), //Stage Use only since its Sync.
		// sdktrace.WithBatcher(exporter), //Production Use, Exports on Flush or Shutdown.
	)
	otel.SetTracerProvider(traceprovider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
}

// ShutdownTracerProvider flushes any pending spans.
// It is recommended to call this function before your program exits
//
// ctx - the context.Context used for shutdown.
// It returns nothing.
func ShutdownTracerProvider(ctx context.Context) {
	traceProvider := otel.GetTracerProvider().(*sdktrace.TracerProvider)
	traceProvider.Shutdown(ctx)
}

// FlushTraceProvider flushes any pending spans.
func FlushTraceProvider(ctx context.Context) {
	traceProvider := otel.GetTracerProvider().(*sdktrace.TracerProvider)
	traceProvider.ForceFlush(ctx)
}
