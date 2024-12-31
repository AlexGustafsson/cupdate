package otelutil

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

// Init initializes Open Telemetry.
// Returns a shutdown function to call before exiting.
func Init(ctx context.Context, target string, insecure bool) (func(context.Context) error, error) {
	// NOTE: otlptracegrpc.New accepts config via OTEL_EXPORTER_OTLP_ENDPOINT and
	// friends. That behavior does not seem to be configurable. For now, let's
	// just keep it undocumented
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(target),
	)
	if err != nil {
		return nil, err
	}

	batchExporter := trace.NewBatchSpanProcessor(exporter)

	resource, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String("cupdate"),
		),
	)
	if err != nil {
		return nil, err
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithResource(resource),
		trace.WithSpanProcessor(batchExporter),
	)

	otel.SetTextMapPropagator(propagation.TraceContext{})
	otel.SetTracerProvider(tracerProvider)

	return tracerProvider.Shutdown, nil
}
