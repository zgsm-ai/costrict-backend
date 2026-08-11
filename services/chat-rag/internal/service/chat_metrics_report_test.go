package service

import (
	"encoding/json"
	"testing"

	"github.com/zgsm-ai/chat-rag/internal/model"
)

func TestBuildResponseMetrics_UpstreamTraceID(t *testing.T) {
	reporter := NewChatMetricsReporter("", "", "X-Custom-Trace-ID")
	chatLog := &model.ChatLog{
		ResponseHeaders: []map[string]string{
			{"x-custom-trace-id": "trace-old"},
			{"X-Custom-Trace-Id": "trace-final"},
		},
	}

	metrics := reporter.buildResponseMetrics(chatLog, nil)
	if metrics.UpstreamTraceID != "trace-final" {
		t.Fatalf("UpstreamTraceID = %q, want %q", metrics.UpstreamTraceID, "trace-final")
	}

	payload, err := json.Marshal(struct {
		ResponseMetrics ResponseMetrics `json:"response_metrics"`
	}{ResponseMetrics: metrics})
	if err != nil {
		t.Fatalf("marshal metrics: %v", err)
	}
	if string(payload) == "" || !containsJSONField(payload, "upstream_trace_id", "trace-final") {
		t.Fatalf("payload does not contain response_metrics.upstream_trace_id: %s", payload)
	}
}

func TestBuildResponseMetrics_UpstreamTraceIDDisabled(t *testing.T) {
	reporter := NewChatMetricsReporter("", "", "")
	chatLog := &model.ChatLog{
		ResponseHeaders: []map[string]string{{"x-oneapi-request-id": "trace-id"}},
	}

	metrics := reporter.buildResponseMetrics(chatLog, nil)
	if metrics.UpstreamTraceID != "" {
		t.Fatalf("UpstreamTraceID = %q, want empty", metrics.UpstreamTraceID)
	}
}

func containsJSONField(payload []byte, name, want string) bool {
	var decoded struct {
		ResponseMetrics map[string]any `json:"response_metrics"`
	}
	if err := json.Unmarshal(payload, &decoded); err != nil {
		return false
	}
	got, ok := decoded.ResponseMetrics[name].(string)
	return ok && got == want
}
