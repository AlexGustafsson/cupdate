package slogutil

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/trace"
)

func Context(ctx context.Context) slog.Attr {
	// TODO: Take traceparent from http requests as well
	span := trace.SpanContextFromContext(ctx)
	if span.HasTraceID() {
		return slog.String("traceId", span.TraceID().String())
	}

	return slog.String("traceId", "")
}
