package slogutil

import (
	"context"
	"io"
	"log/slog"

	"go.opentelemetry.io/otel/trace"
)

// NewHandler initializes the default slog logger to use [Handler].
func NewHandler(w io.Writer, opts *slog.HandlerOptions) slog.Handler {
	return &Handler{
		handler: slog.NewJSONHandler(w, opts),
	}
}

var _ (slog.Handler) = (*Handler)(nil)

// Handler is a [slog.Handler] implementation supporting additions like trace
// ids from context values.
type Handler struct {
	handler slog.Handler
}

// Enabled implements slog.Handler.
func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

// Handle implements slog.Handler.
func (h *Handler) Handle(ctx context.Context, record slog.Record) error {
	span := trace.SpanContextFromContext(ctx)
	if span.HasTraceID() {
		record.AddAttrs(slog.String("traceId", span.TraceID().String()))
	}

	return h.handler.Handle(ctx, record)
}

// WithAttrs implements slog.Handler.
func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Handler{
		handler: h.handler.WithAttrs(attrs),
	}
}

// WithGroup implements slog.Handler.
func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{
		handler: h.handler.WithGroup(name),
	}
}
