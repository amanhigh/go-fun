package common

import (
	"context"

	"github.com/amanhigh/go-fun/common/telemetry"
	"github.com/golobby/container/v3"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	metric_sdk "go.opentelemetry.io/otel/sdk/metric"
)

func (fi *FunAppInjector) setupTelemetry() {
	telemetry.InitLogger(fi.config.Log)
	telemetry.InitTracerProvider(context.Background(), NAMESPACE, fi.config.Tracing)
	setupPrometheus()
}

func setupPrometheus() {
	exporter, err := prometheus.New()
	if err != nil {
		log.Fatal().Err(err).Msg("Prometheus Exporter Failed")
	}

	provider := metric_sdk.NewMeterProvider(
		metric_sdk.WithReader(exporter),
	)
	otel.SetMeterProvider(provider)
}

func (fi *FunAppInjector) registerMetrics() {
	meter := otel.GetMeterProvider().Meter(NAMESPACE)

	container.MustNamedSingleton(fi.di, "CreateCounter", func() metric.Int64Counter {
		counter, _ := meter.Int64Counter("create_person",
			metric.WithDescription("Counts Person Create API"),
		)
		return counter
	})

	container.MustNamedSingleton(fi.di, "PersonCounter", func() metric.Int64UpDownCounter {
		counter, _ := meter.Int64UpDownCounter("person_count",
			metric.WithDescription("Person Count in Get Persons"),
		)
		return counter
	})

	container.MustNamedSingleton(fi.di, "PersonCreateTime", func() metric.Float64Histogram {
		histogram, _ := meter.Float64Histogram("person_create_time",
			metric.WithDescription("Time Taken to Create Person"),
		)
		return histogram
	})
}
