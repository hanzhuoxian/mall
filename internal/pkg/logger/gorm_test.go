package logger_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	gormlogger "gorm.io/gorm/logger"

	gormlog "github.com/hanzhuoxian/mall/internal/pkg/logger"
	"github.com/hanzhuoxian/mall/pkg/logger"
)

func TestGormUsesInjectedLoggerAndRequestContext(t *testing.T) {
	restore := zap.ReplaceGlobals(zap.NewNop())
	t.Cleanup(restore)
	core, logs := observer.New(zap.InfoLevel)
	l := gormlog.NewGormLogger(zap.New(core, zap.AddCaller()), gormlogger.Info)
	sc := trace.NewSpanContext(trace.SpanContextConfig{TraceID: trace.TraceID{1}, SpanID: trace.SpanID{2}})
	ctx := logger.WithFields(trace.ContextWithSpanContext(context.Background(), sc), zap.String(logger.KeyRequestID, "sql-request"))
	l.Info(ctx, "info %s", "detail")
	l.Warn(ctx, "warn %s", "detail")
	l.Error(ctx, "error %s", "detail")
	l.Trace(ctx, time.Now(), func() (string, int64) { return "SELECT 1", 1 }, nil)
	entries := logs.All()
	if len(entries) != 4 {
		t.Fatalf("injected logger received %d logs, want 4", len(entries))
	}
	for _, entry := range entries {
		fields := entry.ContextMap()
		if fields[logger.KeyRequestID] != "sql-request" || fields[logger.KeyTraceID] != sc.TraceID().String() || fields[logger.KeySpanID] != sc.SpanID().String() {
			t.Errorf("%s lost SQL correlation: %+v", entry.Message, fields)
		}
	}
	fields := entries[3].ContextMap()
	if fields["sql"] != "SELECT 1" || fields["rows"] != int64(1) {
		t.Errorf("SQL details missing: %+v", fields)
	}
	if caller, ok := fields["sql_caller"].(string); !ok || !strings.Contains(caller, "gorm_test.go:") {
		t.Errorf("SQL caller missing: %+v", fields)
	}
	if _, duplicate := fields["caller"]; duplicate || !entries[3].Caller.Defined {
		t.Error("SQL caller should be separate from zap's caller")
	}
}

func TestGormTracePoliciesAndLogMode(t *testing.T) {
	core, logs := observer.New(zap.InfoLevel)
	original := gormlog.NewGormLogger(zap.New(core), gormlogger.Info)
	for _, tc := range []struct {
		name    string
		mode    gormlogger.LogLevel
		elapsed time.Duration
		err     error
		want    zapcore.Level
		logged  bool
	}{
		{"normal", gormlogger.Info, 0, nil, zap.InfoLevel, true},
		{"slow", gormlogger.Warn, time.Second, nil, zap.WarnLevel, true},
		{"failure", gormlogger.Error, 0, errors.New("database unavailable"), zap.ErrorLevel, true},
		{"not found", gormlogger.Warn, 0, gormlogger.ErrRecordNotFound, zap.InfoLevel, false},
		{"silent", gormlogger.Silent, time.Second, errors.New("hidden"), zap.InfoLevel, false},
		{"info suppressed", gormlogger.Warn, 0, nil, zap.InfoLevel, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			called := false
			original.LogMode(tc.mode).Trace(context.Background(), time.Now().Add(-tc.elapsed), func() (string, int64) {
				called = true
				return "SELECT 1", 1
			}, tc.err)
			entries := logs.TakeAll()
			if called != tc.logged || (len(entries) == 1) != tc.logged {
				t.Fatalf("logged=%t callback=%t entries=%d", tc.logged, called, len(entries))
			}
			if tc.logged && entries[0].Level != tc.want {
				t.Errorf("level=%v want %v", entries[0].Level, tc.want)
			}
		})
	}
	original.Info(context.Background(), "original remains at info")
	if logs.Len() != 1 {
		t.Fatal("LogMode changed the original logger")
	}
}
