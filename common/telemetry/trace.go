package telemetry

import (
	"context"
	"runtime"

	"github.com/amanhigh/go-fun/models/config"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc"
)

/*
	Resources
	Instruments - https://opentelemetry.io/ecosystem/registry/?language=go&component=instrumentation
	Awesome Telemetry - https://github.com/magsther/awesome-opentelemetry
	OTEL Go - https://opentelemetry.io/docs/instrumentation/go/getting-started/
	Guide - https://medium.com/@emafuma/start-using-opentelemetry-with-go-gin-web-framework-9bebca5abadc
*/

func InitTracerProvider(ctx context.Context, name string, config config.Tracing) {
	subLogger := log.With().Str("Name", name).Str("Type", config.Type).Str("Endpoint", config.Endpoint).Str("Publish", config.Publish).Logger()

	var err error
	switch config.Type {
	case "otlp":
		err = InitOtlpTracerProvider(ctx, name, config)
	case "console":
		err = InitStdoutTracerProvider(config)
	default:
		//No Tracer (Defaults to Global Tracer)
	}

	if err != nil {
		subLogger.Err(err).Msg("Error Initializing Tracer Provider")
	}
	subLogger.Debug().Msg("Tracer Provider Initialized")
}

func InitStdoutTracerProvider(config config.Tracing) (err error) {
	var exporter *stdout.Exporter
	exporter, err = stdout.New(stdout.WithPrettyPrint())
	traceprovider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		getPublisher(config, exporter),
	)
	otel.SetTracerProvider(traceprovider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return
}

// https://observiq.com/blog/tracing-services-using-otel-and-jaeger
func InitOtlpTracerProvider(ctx context.Context, name string, config config.Tracing) (err error) {
	var conn *grpc.ClientConn
	var exporter *otlptrace.Exporter

	if conn, err = grpc.DialContext(ctx, config.Endpoint, grpc.WithInsecure()); err == nil {
		if exporter, err = otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure(), otlptracegrpc.WithGRPCConn(conn)); err == nil {
			if resources, err := buildResource(name); err == nil {
				traceprovider := sdktrace.NewTracerProvider(
					//https://opentelemetry.io/docs/instrumentation/go/sampling/
					sdktrace.WithSampler(sdktrace.AlwaysSample()),
					sdktrace.WithResource(resources),
					getPublisher(config, exporter),
				)
				otel.SetTracerProvider(traceprovider)
				otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
			}
		}
	}
	return
}

func getPublisher(config config.Tracing, exporter sdktrace.SpanExporter) trace.TracerProviderOption {
	var publisher sdktrace.TracerProviderOption
	switch config.Publish {
	case "sync":
		//Stage Use only since its Sync.
		publisher = sdktrace.WithSyncer(exporter)
	case "batch":
		//Production Use, Exports on Flush or Shutdown.
		publisher = sdktrace.WithBatcher(exporter)
	}
	return publisher
}

// https://opentelemetry.io/docs/instrumentation/go/resources/
func buildResource(name string) (resources *resource.Resource, err error) {
	resources, err = resource.New(context.Background(),
		resource.WithProcess(),   // This option configures a set of Detectors that discover process information
		resource.WithOS(),        // This option configures a set of Detectors that discover OS information
		resource.WithContainer(), // This option configures a set of Detectors that discover container information
		resource.WithHost(),      // This option configures a set of Detectors that discover host information
		resource.WithAttributes( // Or specify resource attributes directly
			attribute.String("foo", "bar"),
			semconv.ServiceNameKey.String(name),
			semconv.HostArchKey.String(runtime.GOARCH),
		),
	)
	return
}

// ShutdownTracerProvider flushes any pending spans.
// It is recommended to call this function before your program exits
//
// ctx - the context.Context used for shutdown.
// It returns nothing.
func ShutdownTracerProvider(ctx context.Context) {
	if traceProvider, ok := otel.GetTracerProvider().(*sdktrace.TracerProvider); ok {
		traceProvider.Shutdown(ctx)
	}
}

// FlushTraceProvider flushes any pending spans.
func FlushTraceProvider(ctx context.Context) {
	if traceProvider := otel.GetTracerProvider().(*sdktrace.TracerProvider); traceProvider != nil {
		traceProvider.ForceFlush(ctx)
	}
}
