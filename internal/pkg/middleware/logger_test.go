package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"

	"github.com/hanzhuoxian/mall/internal/pkg/middleware"
	"github.com/hanzhuoxian/mall/internal/pkg/response"
	"github.com/hanzhuoxian/mall/pkg/logger"
)

func TestRequestContextReachesHandlerResponseAndAccessLogs(t *testing.T) {
	core, logs := observer.New(zap.InfoLevel)
	restore := zap.ReplaceGlobals(zap.New(core))
	t.Cleanup(restore)
	provider := sdktrace.NewTracerProvider()
	t.Cleanup(func() { _ = provider.Shutdown(context.Background()) })
	engine := gin.New()
	// Gin's context fallback is deliberately disabled (the default).
	engine.Use(otelgin.Middleware("mall-test", otelgin.WithTracerProvider(provider)))
	engine.Use(requestid.New(), middleware.Context())
	engine.Use(middleware.LoggerWithConfig(gin.LoggerConfig{
		Formatter: func(gin.LogFormatterParams) string { return "access" },
	}))
	var requestSpan trace.SpanContext
	engine.GET("/test", func(c *gin.Context) {
		requestSpan = trace.SpanContextFromContext(c.Request.Context())
		logger.L(c.Request.Context()).Info("handler")
		response.Fail(c, &response.HTTPError{Status: http.StatusConflict, Message: "conflict"})
	})
	var previousTrace trace.TraceID
	for _, id := range []string{"request-one", "request-two"} {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("X-Request-ID", id)
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)
		if w.Code != http.StatusConflict || w.Header().Get("X-Request-ID") != id {
			t.Fatalf("unexpected response: %d %+v", w.Code, w.Header())
		}
		if !requestSpan.IsValid() || requestSpan.TraceID() == previousTrace {
			t.Fatalf("missing or reused request span: %v", requestSpan)
		}
		previousTrace = requestSpan.TraceID()
		entries := logs.TakeAll()
		if len(entries) != 3 {
			t.Fatalf("want handler, error and access logs, got %d", len(entries))
		}
		for _, entry := range entries {
			fields := entry.ContextMap()
			if fields[logger.KeyRequestID] != id || fields[logger.KeyTraceID] != requestSpan.TraceID().String() || fields[logger.KeySpanID] != requestSpan.SpanID().String() {
				t.Errorf("%s lost request correlation: %+v", entry.Message, fields)
			}
		}
	}
}
