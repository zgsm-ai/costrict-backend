package service

import (
	"testing"
	"time"
)

func TestErrorLogSampler_AllowsUnderLimit(t *testing.T) {
	s := NewErrorLogSampler(3, 60)
	for i := 0; i < 3; i++ {
		if !s.Allow("alice", "ApiError") {
			t.Fatalf("call %d: expected allow", i+1)
		}
	}
}

func TestErrorLogSampler_DeniesOverLimit(t *testing.T) {
	s := NewErrorLogSampler(1, 60)
	if !s.Allow("alice", "ApiError") {
		t.Fatal("expected first allow")
	}
	if s.Allow("alice", "ApiError") {
		t.Fatal("expected deny after limit reached")
	}
}

func TestErrorLogSampler_ResetsAfterWindow(t *testing.T) {
	s := NewErrorLogSampler(1, 60)
	base := time.Unix(1_000_000_020, 0) // aligned to a 60s window boundary
	s.now = func() time.Time { return base }

	if !s.Allow("alice", "ApiError") {
		t.Fatal("expected first allow")
	}
	if s.Allow("alice", "ApiError") {
		t.Fatal("expected deny within window")
	}

	// Still inside the same window at +59s.
	s.now = func() time.Time { return base.Add(59 * time.Second) }
	if s.Allow("alice", "ApiError") {
		t.Fatal("expected deny at +59s (same window)")
	}

	// New window at +60s.
	s.now = func() time.Time { return base.Add(60 * time.Second) }
	if !s.Allow("alice", "ApiError") {
		t.Fatal("expected allow after window reset at +60s")
	}
}

func TestErrorLogSampler_UsersAreIsolated(t *testing.T) {
	s := NewErrorLogSampler(1, 60)
	if !s.Allow("alice", "ApiError") {
		t.Fatal("expected allow for alice")
	}
	if !s.Allow("bob", "ApiError") {
		t.Fatal("expected allow for bob (separate user bucket)")
	}
	if s.Allow("alice", "ApiError") {
		t.Fatal("expected alice bucket to be exhausted")
	}
}

func TestErrorLogSampler_TypesAreIsolated(t *testing.T) {
	s := NewErrorLogSampler(1, 60)
	if !s.Allow("alice", "ApiError") {
		t.Fatal("expected allow for ApiError")
	}
	if !s.Allow("alice", "ServerError") {
		t.Fatal("expected allow for ServerError (separate type bucket)")
	}
	if s.Allow("alice", "ApiError") {
		t.Fatal("expected ApiError bucket to be exhausted")
	}
}

func TestErrorLogSampler_NonPositiveWindowDoesNotPanic(t *testing.T) {
	s := NewErrorLogSampler(1, 0)
	if !s.Allow("alice", "ApiError") {
		t.Fatal("expected allow with defensive window")
	}
}
