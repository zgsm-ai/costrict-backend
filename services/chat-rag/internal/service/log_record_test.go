package service

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/zgsm-ai/chat-rag/internal/config"
	"github.com/zgsm-ai/chat-rag/internal/model"
	"github.com/zgsm-ai/chat-rag/internal/storage"
	"github.com/zgsm-ai/chat-rag/internal/types"
)

func newErrorLog(user string) *model.ChatLog {
	log := &model.ChatLog{
		Error: []map[types.ErrorType]string{
			{types.ErrApiError: "boom"},
		},
	}
	log.Identity.UserName = user
	return log
}

func TestShouldSaveErrorLog_AllMode(t *testing.T) {
	ls := &LoggerRecordService{errorLogMode: config.ErrorLogModeAll}
	if !ls.shouldSaveErrorLog(newErrorLog("alice")) {
		t.Fatal("expected save in all mode")
	}
}

func TestShouldSaveErrorLog_NoneMode(t *testing.T) {
	ls := &LoggerRecordService{errorLogMode: config.ErrorLogModeNone}
	if ls.shouldSaveErrorLog(newErrorLog("alice")) {
		t.Fatal("expected skip in none mode")
	}
}

func TestShouldSaveErrorLog_SampledMode(t *testing.T) {
	ls := &LoggerRecordService{
		errorLogMode: config.ErrorLogModeSampled,
		errorSampler: NewErrorLogSampler(1, 60),
	}
	if !ls.shouldSaveErrorLog(newErrorLog("alice")) {
		t.Fatal("expected first sampled save for alice to be allowed")
	}
	if ls.shouldSaveErrorLog(newErrorLog("alice")) {
		t.Fatal("expected second sampled save for alice to be dropped")
	}
	if !ls.shouldSaveErrorLog(newErrorLog("bob")) {
		t.Fatal("expected sampled save for bob (separate user) to be allowed")
	}
}

func TestShouldSaveErrorLog_SampledModeNilSampler(t *testing.T) {
	ls := &LoggerRecordService{errorLogMode: config.ErrorLogModeSampled}
	if !ls.shouldSaveErrorLog(newErrorLog("alice")) {
		t.Fatal("expected save when sampler is nil (fail-open)")
	}
}

func TestFirstErrorTypeKey(t *testing.T) {
	if got := firstErrorTypeKey(newErrorLog("alice")); got != string(types.ErrApiError) {
		t.Errorf("firstErrorTypeKey() = %q, want %q", got, types.ErrApiError)
	}
	if got := firstErrorTypeKey(&model.ChatLog{}); got != "" {
		t.Errorf("firstErrorTypeKey() on empty = %q, want empty", got)
	}
	emptyMap := &model.ChatLog{Error: []map[types.ErrorType]string{{}}}
	if got := firstErrorTypeKey(emptyMap); got != "" {
		t.Errorf("firstErrorTypeKey() on empty map = %q, want empty", got)
	}
}

// --- Integration tests: full wiring from config.Config → NewLogRecordService → shouldSaveErrorLog ---

func TestNewLogRecordService_Integration_ErrorLogModeAll(t *testing.T) {
	cfg := config.Config{}
	cfg.Log.ErrorLogMode = config.ErrorLogModeAll

	svc := NewLogRecordService(cfg)
	ls, ok := svc.(*LoggerRecordService)
	if !ok {
		t.Fatal("NewLogRecordService did not return *LoggerRecordService")
	}

	if ls.errorLogMode != config.ErrorLogModeAll {
		t.Errorf("errorLogMode = %q, want %q", ls.errorLogMode, config.ErrorLogModeAll)
	}
	if ls.errorSampler != nil {
		t.Error("errorSampler should be nil in all mode")
	}
	if !ls.shouldSaveErrorLog(newErrorLog("alice")) {
		t.Fatal("expected save in all mode via integration path")
	}
}

func TestNewLogRecordService_Integration_ErrorLogModeNone(t *testing.T) {
	cfg := config.Config{}
	cfg.Log.ErrorLogMode = config.ErrorLogModeNone

	svc := NewLogRecordService(cfg)
	ls, ok := svc.(*LoggerRecordService)
	if !ok {
		t.Fatal("NewLogRecordService did not return *LoggerRecordService")
	}

	if ls.errorLogMode != config.ErrorLogModeNone {
		t.Errorf("errorLogMode = %q, want %q", ls.errorLogMode, config.ErrorLogModeNone)
	}
	if ls.errorSampler != nil {
		t.Error("errorSampler should be nil in none mode")
	}
	if ls.shouldSaveErrorLog(newErrorLog("alice")) {
		t.Fatal("expected skip in none mode via integration path")
	}
}

func TestNewLogRecordService_Integration_ErrorLogModeSampled(t *testing.T) {
	cfg := config.Config{}
	cfg.Log.ErrorLogMode = config.ErrorLogModeSampled

	svc := NewLogRecordService(cfg)
	ls, ok := svc.(*LoggerRecordService)
	if !ok {
		t.Fatal("NewLogRecordService did not return *LoggerRecordService")
	}

	if ls.errorLogMode != config.ErrorLogModeSampled {
		t.Errorf("errorLogMode = %q, want %q", ls.errorLogMode, config.ErrorLogModeSampled)
	}
	if ls.errorSampler == nil {
		t.Fatal("errorSampler should not be nil in sampled mode")
	}

	// Verify the sampler actually limits: first call allowed, second denied.
	if !ls.shouldSaveErrorLog(newErrorLog("alice")) {
		t.Fatal("expected first sampled save for alice to be allowed")
	}
	if ls.shouldSaveErrorLog(newErrorLog("alice")) {
		t.Fatal("expected second sampled save for alice to be dropped")
	}
	// Different user should still be allowed.
	if !ls.shouldSaveErrorLog(newErrorLog("bob")) {
		t.Fatal("expected sampled save for bob (separate user) to be allowed")
	}
}

func TestNewLogRecordService_Integration_SamplerDefaultsApplied(t *testing.T) {
	cfg := config.Config{}
	cfg.Log.ErrorLogMode = config.ErrorLogModeSampled
	// Leave ErrorLogSampleN and ErrorLogSampleWindowSec at zero to exercise defaults.

	svc := NewLogRecordService(cfg)
	ls, ok := svc.(*LoggerRecordService)
	if !ok {
		t.Fatal("NewLogRecordService did not return *LoggerRecordService")
	}

	if ls.errorSampler == nil {
		t.Fatal("errorSampler should not be nil when defaults are applied")
	}
	if ls.errorSampler.maxPerWindow != 1 {
		t.Errorf("default maxPerWindow = %d, want 1", ls.errorSampler.maxPerWindow)
	}
	if ls.errorSampler.windowSec != 60 {
		t.Errorf("default windowSec = %d, want 60", ls.errorSampler.windowSec)
	}
}

func TestNewLogRecordService_Integration_EmptyUserNameFallback(t *testing.T) {
	cfg := config.Config{}
	cfg.Log.ErrorLogMode = config.ErrorLogModeSampled

	svc := NewLogRecordService(cfg)
	ls, ok := svc.(*LoggerRecordService)
	if !ok {
		t.Fatal("NewLogRecordService did not return *LoggerRecordService")
	}

	// Log with empty UserName: shouldSaveErrorLog must use "unknown" fallback
	// via sanitizeFilename, and the first call must be allowed.
	log := newErrorLog("")
	if !ls.shouldSaveErrorLog(log) {
		t.Fatal("expected save for empty username (fallback to 'unknown')")
	}
	// Second call with same empty user should be denied (N=1 window).
	if ls.shouldSaveErrorLog(log) {
		t.Fatal("expected second empty-username call to be dropped")
	}
	// A different error type for the same empty user should still be allowed
	// (separate type bucket).
	logDiffType := &model.ChatLog{
		Error: []map[types.ErrorType]string{
			{types.ErrServerError: "internal error"},
		},
	}
	logDiffType.Identity.UserName = ""
	if !ls.shouldSaveErrorLog(logDiffType) {
		t.Fatal("expected save for empty username with different error type")
	}
}

func TestLoggerRecordService_StartDoesNotCreateDirectoryWhenStorageDisabled(t *testing.T) {
	logPath := filepath.Join(t.TempDir(), "disabled-logs")
	cfg := config.Config{}
	cfg.Log.LogFilePath = logPath
	cfg.Log.StorageType = "none"

	svc := NewLogRecordService(cfg)
	ls, ok := svc.(*LoggerRecordService)
	if !ok {
		t.Fatal("NewLogRecordService did not return *LoggerRecordService")
	}
	if err := ls.Start(); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	ls.Stop()

	if _, err := os.Stat(logPath); !os.IsNotExist(err) {
		t.Fatalf("expected log directory not to exist, stat error = %v", err)
	}
}

// --- Regression tests ---

func TestShouldSaveErrorLog_UnknownMode(t *testing.T) {
	// When errorLogMode is set to an unrecognized value (neither "all",
	// "sampled", nor "none"), the default case in shouldSaveErrorLog must
	// deny the save — fail closed, not open.
	ls := &LoggerRecordService{errorLogMode: "bogus"}
	if ls.shouldSaveErrorLog(newErrorLog("alice")) {
		t.Fatal("expected deny for unknown errorLogMode (fail-closed)")
	}
}

func TestShouldSaveErrorLog_DifferentErrorTypesIndependentBuckets(t *testing.T) {
	// Per-user-per-type sampling: different error types for the same user
	// are tracked in separate buckets, so each type gets its own allowance.
	ls := &LoggerRecordService{
		errorLogMode: config.ErrorLogModeSampled,
		errorSampler: NewErrorLogSampler(1, 60),
	}

	logApiErr := &model.ChatLog{
		Error: []map[types.ErrorType]string{
			{types.ErrApiError: "api failure"},
		},
	}
	logApiErr.Identity.UserName = "alice"

	logServerErr := &model.ChatLog{
		Error: []map[types.ErrorType]string{
			{types.ErrServerError: "internal failure"},
		},
	}
	logServerErr.Identity.UserName = "alice"

	// First call with ApiError: allowed.
	if !ls.shouldSaveErrorLog(logApiErr) {
		t.Fatal("expected first ApiError save for alice to be allowed")
	}
	// First call with ServerError (different type): still allowed.
	if !ls.shouldSaveErrorLog(logServerErr) {
		t.Fatal("expected first ServerError save for alice to be allowed")
	}
	// Second call with ApiError: denied (same bucket exhausted).
	if ls.shouldSaveErrorLog(logApiErr) {
		t.Fatal("expected second ApiError save for alice to be dropped")
	}
}

// --- Invariant: metrics are reported even when the error log is not persisted ---

type fakeMetrics struct{ recordCalls int }

func (f *fakeMetrics) RecordChatLog(*model.ChatLog)      { f.recordCalls++ }
func (f *fakeMetrics) GetRegistry() *prometheus.Registry { return nil }

type fakeStorage struct{ writes int }

func (f *fakeStorage) Write(key string, data []byte) (*storage.WriteInfo, error) {
	f.writes++
	return &storage.WriteInfo{FilePath: key}, nil
}
func (f *fakeStorage) Close() error { return nil }

func newErrorLogWithUser(user string) *model.ChatLog {
	log := newErrorLog(user)
	log.Identity.UserInfo = &model.UserInfo{}
	return log
}

func TestLogDirectToStorage_NoneModeStillReportsMetrics(t *testing.T) {
	metrics := &fakeMetrics{}
	store := &fakeStorage{}
	ls := &LoggerRecordService{errorLogMode: config.ErrorLogModeNone}
	ls.SetMetricsService(metrics)
	ls.SetStorageBackend(store)

	ls.logDirectToStorage(newErrorLogWithUser("alice"))

	if metrics.recordCalls != 1 {
		t.Fatalf("expected RecordChatLog called once even when save skipped, got %d", metrics.recordCalls)
	}
	if store.writes != 0 {
		t.Fatalf("expected no storage write in none mode, got %d", store.writes)
	}
}

func TestLogDirectToStorage_SampledReportsMetricsOnDrop(t *testing.T) {
	metrics := &fakeMetrics{}
	store := &fakeStorage{}
	ls := &LoggerRecordService{
		errorLogMode: config.ErrorLogModeSampled,
		errorSampler: NewErrorLogSampler(1, 60),
	}
	ls.SetMetricsService(metrics)
	ls.SetStorageBackend(store)

	for i := 0; i < 2; i++ {
		ls.logDirectToStorage(newErrorLogWithUser("alice"))
	}

	if metrics.recordCalls != 2 {
		t.Fatalf("expected RecordChatLog called twice (always reported), got %d", metrics.recordCalls)
	}
	if store.writes != 1 {
		t.Fatalf("expected exactly one storage write (second sampled out), got %d", store.writes)
	}
}

func TestLogDirectToStorage_DisabledStillReportsMetricsAndChatMetrics(t *testing.T) {
	reported := make(chan struct{}, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reported <- struct{}{}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	metrics := &fakeMetrics{}
	store := &fakeStorage{}
	ls := &LoggerRecordService{
		storageDisabled: true,
		metricsReporter: NewChatMetricsReporter(server.URL, http.MethodPost, ""),
	}
	ls.SetMetricsService(metrics)
	ls.SetStorageBackend(store)

	chatLog := &model.ChatLog{}
	chatLog.Identity.RequestID = "storage-disabled-request"
	chatLog.Identity.UserInfo = &model.UserInfo{}
	ls.logDirectToStorage(chatLog)

	if metrics.recordCalls != 1 {
		t.Fatalf("expected RecordChatLog called once, got %d", metrics.recordCalls)
	}
	if store.writes != 0 {
		t.Fatalf("expected no storage writes when storage is disabled, got %d", store.writes)
	}

	select {
	case <-reported:
	case <-time.After(2 * time.Second):
		t.Fatal("expected chat metrics report when storage is disabled")
	}
}
