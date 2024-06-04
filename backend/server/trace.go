package server

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/trace"
)

type traceLoggingHandler struct {
	inner slog.Handler
}

var _ slog.Handler = (*traceLoggingHandler)(nil)

func NewTraceLoggingHandler(inner slog.Handler) *traceLoggingHandler {
	return &traceLoggingHandler{inner}
}

func (h *traceLoggingHandler) Handle(ctx context.Context, r slog.Record) error {
	sc := trace.SpanContextFromContext(ctx)
	if sc.IsValid() {
		r.AddAttrs(
			slog.String("trace_id", sc.TraceID().String()),
		)
	}
	return h.inner.Handle(ctx, r)
}

func (h *traceLoggingHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.inner.Enabled(ctx, level)
}

func (h *traceLoggingHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &traceLoggingHandler{h.inner.WithAttrs(attrs)}
}

func (h *traceLoggingHandler) WithGroup(name string) slog.Handler {
	return &traceLoggingHandler{h.inner.WithGroup(name)}
}
