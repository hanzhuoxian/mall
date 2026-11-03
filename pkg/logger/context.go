package logger

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Common structured log field names.
const (
	KeyRequestID = "requestID"
	KeyUsername  = "username"
	KeyTraceID   = "trace_id"
	KeySpanID    = "span_id"
)

type fieldsKey struct{}

// WithFields adds request metadata to a standard context, preserving parent fields.
func WithFields(ctx context.Context, fields ...zap.Field) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	parent, _ := ctx.Value(fieldsKey{}).([]zap.Field)
	combined := make([]zap.Field, 0, len(parent)+len(fields))
	combined = append(combined, parent...)
	combined = append(combined, fields...)
	return context.WithValue(ctx, fieldsKey{}, combined)
}

// L returns the global zap logger enriched with request metadata and the current span.
func L(ctx context.Context) *zap.Logger {
	return WithContext(ctx, zap.L())
}

// WithContext enriches an injected logger. Context is passed to otelzap as a skip
// field for native OTLP trace correlation, and IDs are also included in local output.
func WithContext(ctx context.Context, l *zap.Logger) *zap.Logger {
	if ctx == nil {
		ctx = context.Background()
	}
	stored, _ := ctx.Value(fieldsKey{}).([]zap.Field)
	fields := make([]zap.Field, 0, len(stored)+3)
	fields = append(fields, stored...)
	fields = append(fields, zap.Field{Type: zapcore.SkipType, Interface: ctx})
	if sc := trace.SpanContextFromContext(ctx); sc.IsValid() {
		fields = append(fields,
			zap.String(KeyTraceID, sc.TraceID().String()),
			zap.String(KeySpanID, sc.SpanID().String()))
	}
	return l.With(fields...)
}
