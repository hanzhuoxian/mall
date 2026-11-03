package logger_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	otelapi "go.opentelemetry.io/otel/log"
	otellog "go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/hanzhuoxian/mall/pkg/logger"
)

// recordSink observes records after real zap -> otelzap -> SDK processing.
type recordSink struct {
	mu      sync.Mutex
	records []sdklog.Record
}

func (s *recordSink) OnEmit(_ context.Context, r *sdklog.Record) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.records = append(s.records, r.Clone())
	return nil
}
func (*recordSink) Shutdown(context.Context) error                         { return nil }
func (*recordSink) ForceFlush(context.Context) error                       { return nil }
func (*recordSink) Enabled(context.Context, sdklog.EnabledParameters) bool { return true }
func (s *recordSink) Records() []sdklog.Record {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]sdklog.Record(nil), s.records...)
}

func setup(t *testing.T, level string) (*logger.Options, string, *recordSink) {
	t.Helper()
	previousLogger, previousProvider := zap.L(), otellog.GetLoggerProvider()
	previousOutput, previousFlags, previousPrefix := log.Writer(), log.Flags(), log.Prefix()
	sink := &recordSink{}
	provider := sdklog.NewLoggerProvider(sdklog.WithProcessor(sink))
	otellog.SetLoggerProvider(provider)
	t.Cleanup(func() {
		_ = zap.L().Sync()
		zap.ReplaceGlobals(previousLogger)
		log.SetOutput(previousOutput)
		log.SetFlags(previousFlags)
		log.SetPrefix(previousPrefix)
		otellog.SetLoggerProvider(previousProvider)
		_ = provider.Shutdown(context.Background())
	})
	path := filepath.Join(t.TempDir(), "app.log")
	opts := logger.NewOptions()
	opts.Format, opts.Level, opts.Name = "json", level, "mall-test"
	opts.OutputPaths = []string{path}
	return opts, path, sink
}

func readLogs(t *testing.T, path string) []map[string]any {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var entries []map[string]any
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		if line == "" {
			continue
		}
		var entry map[string]any
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			t.Fatal(err)
		}
		entries = append(entries, entry)
	}
	return entries
}

func TestInitConnectsNativeZapAndStandardLog(t *testing.T) {
	opts, path, sink := setup(t, "info")
	if err := logger.Init(opts); err != nil {
		t.Fatal(err)
	}
	zap.L().Info("native")
	zap.S().Infof("sugared %s", "message")
	log.Print("standard")
	entries := readLogs(t, path)
	if len(entries) != 3 || len(sink.Records()) != 3 {
		t.Fatalf("native/sugared/standard logs: local=%d otel=%d, want 3 each", len(entries), len(sink.Records()))
	}
}

func TestLevelFiltersBothSinks(t *testing.T) {
	opts, path, sink := setup(t, "warn")
	if err := logger.Init(opts); err != nil {
		t.Fatal(err)
	}
	l := logger.L(context.Background())
	l.Debug("hidden debug")
	l.Info("hidden info")
	l.Warn("visible warn")
	l.Error("visible error")
	entries := readLogs(t, path)
	if len(entries) != 2 || len(sink.Records()) != 2 {
		t.Fatalf("warn threshold: local=%d otel=%d, want 2 each", len(entries), len(sink.Records()))
	}
}

func TestContextCorrelatesLocalAndOTLPLogs(t *testing.T) {
	opts, path, sink := setup(t, "info")
	if err := logger.Init(opts); err != nil {
		t.Fatal(err)
	}
	sc := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: trace.TraceID{1}, SpanID: trace.SpanID{2}, TraceFlags: trace.FlagsSampled,
	})
	ctx := trace.ContextWithSpanContext(context.Background(), sc)
	logger.L(ctx).Info("request")
	entries, records := readLogs(t, path), sink.Records()
	if len(entries) != 1 || entries[0]["trace_id"] != sc.TraceID().String() || entries[0]["span_id"] != sc.SpanID().String() {
		t.Fatalf("local correlation missing: %+v", entries)
	}
	if len(records) != 1 || records[0].TraceID() != sc.TraceID() || records[0].SpanID() != sc.SpanID() {
		t.Fatalf("OTLP correlation missing: %+v", records)
	}
}

// A fixed entry timestamp makes the sampler's one-second window deterministic.
type fixedClock struct{ time.Time }

func (c fixedClock) Now() time.Time                       { return c.Time }
func (fixedClock) NewTicker(d time.Duration) *time.Ticker { return time.NewTicker(d) }

func TestSamplingMatchesAcrossSinks(t *testing.T) {
	opts, path, sink := setup(t, "info")
	l, err := logger.New(opts)
	if err != nil {
		t.Fatal(err)
	}
	l = l.WithOptions(zap.WithClock(fixedClock{time.Now()}))
	for i := 1; i <= 250; i++ {
		l.Info("repeated", zap.Int("sequence", i))
	}
	entries, records := readLogs(t, path), sink.Records()
	if len(entries) != 101 || len(records) != len(entries) {
		t.Fatalf("sampling: local=%d otel=%d, want 101 each", len(entries), len(records))
	}
	for i, record := range records {
		var sequence int64
		record.WalkAttributes(func(attr otelapi.KeyValue) bool {
			if attr.Key == "sequence" {
				sequence = attr.Value.AsInt64()
			}
			return true
		})
		if entries[i]["sequence"] != float64(sequence) {
			t.Fatalf("sinks sampled different entries: local=%v otel=%d", entries[i], sequence)
		}
	}
}

func TestNewAndInvalidInitPreserveGlobals(t *testing.T) {
	opts, _, _ := setup(t, "info")
	before, output := zap.L(), log.Writer()
	opts.Format = "JSON"
	l, err := logger.New(opts)
	if err != nil {
		t.Fatal(err)
	}
	defer l.Sync()
	if zap.L() != before || log.Writer() != output {
		t.Fatal("New changed global logging")
	}
	for _, tc := range []struct {
		name string
		edit func(*logger.Options)
	}{
		{"level", func(o *logger.Options) { o.Level = "invalid" }},
		{"format", func(o *logger.Options) { o.Format = "invalid" }},
		{"output", func(o *logger.Options) { o.OutputPaths = []string{filepath.Join(t.TempDir(), "missing", "app.log")} }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			invalid := *opts
			tc.edit(&invalid)
			if err := logger.Init(&invalid); err == nil {
				t.Fatal("invalid configuration accepted")
			}
			if zap.L() != before || log.Writer() != output {
				t.Fatal("failed Init replaced global logging")
			}
		})
	}
}

func TestCallerAndStacktraceOptions(t *testing.T) {
	for _, disabled := range []bool{false, true} {
		t.Run(fmt.Sprintf("disabled=%t", disabled), func(t *testing.T) {
			opts, path, _ := setup(t, "info")
			opts.DisableStacktrace = disabled
			if err := logger.Init(opts); err != nil {
				t.Fatal(err)
			}
			_, _, nativeLine, _ := runtime.Caller(0)
			logger.L(nil).Info("native caller")
			_, _, sugarLine, _ := runtime.Caller(0)
			logger.L(nil).Sugar().Infof("sugar %s", "caller")
			func() {
				defer func() {
					if recover() == nil {
						t.Error("Panic did not panic")
					}
				}()
				logger.L(nil).Panic("panic stack")
			}()
			entries := readLogs(t, path)
			if len(entries) != 3 {
				t.Fatalf("want 3 records, got %d", len(entries))
			}
			for i, line := range []int{nativeLine + 1, sugarLine + 1} {
				if entries[i]["caller"] != fmt.Sprintf("logger/log_test.go:%d", line) {
					t.Errorf("caller=%v, want log_test.go:%d", entries[i]["caller"], line)
				}
			}
			_, hasStack := entries[2]["stacktrace"]
			if hasStack == disabled {
				t.Errorf("DisableStacktrace=%t, stacktrace present=%t", disabled, hasStack)
			}
		})
	}
}

func TestContextFieldsDoNotLeakBetweenRequests(t *testing.T) {
	opts, path, sink := setup(t, "info")
	if err := logger.Init(opts); err != nil {
		t.Fatal(err)
	}
	parent := logger.WithFields(nil, zap.String("service", "mall"))
	first := logger.WithFields(parent, zap.String(logger.KeyRequestID, "first"))
	second := logger.WithFields(parent, zap.String(logger.KeyRequestID, "second"))
	firstSpan := trace.NewSpanContext(trace.SpanContextConfig{TraceID: trace.TraceID{1}, SpanID: trace.SpanID{1}})
	childSpan := trace.NewSpanContext(trace.SpanContextConfig{TraceID: firstSpan.TraceID(), SpanID: trace.SpanID{2}})
	first = trace.ContextWithSpanContext(first, firstSpan)
	logger.L(first).Info("first")
	logger.L(trace.ContextWithSpanContext(first, childSpan)).Info("child")
	logger.L(second).Info("second")
	logger.L(parent).Info("parent")
	entries, records := readLogs(t, path), sink.Records()
	if len(entries) != 4 || len(records) != 4 {
		t.Fatalf("want 4 records, local=%d otel=%d", len(entries), len(records))
	}
	for i, want := range []any{"first", "first", "second", nil} {
		if entries[i][logger.KeyRequestID] != want || entries[i]["service"] != "mall" {
			t.Errorf("context fields at %d: %+v", i, entries[i])
		}
	}
	if records[0].SpanID() != firstSpan.SpanID() || records[1].SpanID() != childSpan.SpanID() || records[2].SpanID().IsValid() {
		t.Fatal("span context was lost, stale or leaked")
	}
}

func TestConcurrentInitAndLogging(t *testing.T) {
	opts, _, sink := setup(t, "info")
	opts.OutputPaths = nil
	if err := logger.Init(opts); err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	for worker := 0; worker < 4; worker++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 10; i++ {
				if err := logger.Init(opts); err != nil {
					t.Error(err)
					return
				}
				zap.L().Info("native concurrent")
				logger.L(context.Background()).Info("context concurrent")
				log.Print("standard concurrent")
			}
		}()
	}
	wg.Wait()
	if got := len(sink.Records()); got != 120 {
		t.Fatalf("lost concurrent logs: got %d, want 120", got)
	}
}
