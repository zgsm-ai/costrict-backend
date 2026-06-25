package config

import "testing"

func TestResolveErrorLogMode(t *testing.T) {
	tests := []struct {
		name         string
		mode         string
		saveErrorLog bool
		want         string
	}{
		{"explicit all", "all", false, "all"},
		{"explicit sampled", "sampled", false, "sampled"},
		{"explicit none", "none", true, "none"},
		{"trims and lowercases", "  Sampled ", false, "sampled"},
		{"uppercase", "ALL", false, "all"},
		{"empty falls back to true->all", "", true, "all"},
		{"empty falls back to false->none", "", false, "none"},
		{"invalid falls back to true->all", "bogus", true, "all"},
		{"invalid falls back to false->none", "bogus", false, "none"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := LogConfig{ErrorLogMode: tt.mode, SaveErrorLog: tt.saveErrorLog}
			if got := c.ResolveErrorLogMode(); got != tt.want {
				t.Errorf("ResolveErrorLogMode() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestErrorLogSamplingEnabled(t *testing.T) {
	if !(LogConfig{ErrorLogMode: "sampled"}).ErrorLogSamplingEnabled() {
		t.Error("expected sampling enabled for sampled mode")
	}
	if (LogConfig{ErrorLogMode: "all"}).ErrorLogSamplingEnabled() {
		t.Error("expected sampling disabled for all mode")
	}
	if (LogConfig{ErrorLogMode: "", SaveErrorLog: true}).ErrorLogSamplingEnabled() {
		t.Error("expected sampling disabled for fallback all")
	}
}
