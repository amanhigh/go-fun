package metrics

import (
	"context"
	"os"
	"runtime"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc"
)

/*
	Resources
	Instruments - https://opentelemetry.io/ecosystem/registry/?language=go&component=instrumentation
	Awesome Telemetry - https://github.com/magsther/awesome-opentelemetry
	OTEL Go - https://opentelemetry.io/docs/instrumentation/go/getting-started/
*/

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
