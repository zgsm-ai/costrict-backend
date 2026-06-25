package service

import (
	"sync"
	"time"
)

// ErrorLogSampler bounds how many error logs are persisted per (user, error type)
// within each fixed time window. It is safe for concurrent use. Counters live in
// memory and are per-instance. The window is aligned globally as
// unixSeconds / windowSec; when the window number changes the whole counter map
// is dropped, which both aligns windows and bounds memory (no cross-window growth).
type ErrorLogSampler struct {
	maxPerWindow int
	windowSec    int64
	now          func() time.Time

	mu       sync.Mutex
	windowID int64
	counts   map[string]int
}

// NewErrorLogSampler creates a sampler that keeps at most maxPerWindow entries
// per (user, error type) within each window of windowSec seconds. A non-positive
// windowSec is clamped to 1 second.
func NewErrorLogSampler(maxPerWindow int, windowSec int64) *ErrorLogSampler {
	if windowSec <= 0 {
		windowSec = 1
	}
	if maxPerWindow <= 0 {
		maxPerWindow = 1
	}
	return &ErrorLogSampler{
		maxPerWindow: maxPerWindow,
		windowSec:    windowSec,
		now:          time.Now,
		windowID:     -1,
		counts:       make(map[string]int),
	}
}

// Allow reports whether an error log for the given user and error type may be
// saved in the current window. When it returns true, the call is counted.
func (s *ErrorLogSampler) Allow(user, errorType string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	wid := s.now().Unix() / s.windowSec
	if wid != s.windowID {
		s.windowID = wid
		s.counts = make(map[string]int)
	}

	key := user + "\x00" + errorType
	if s.counts[key] >= s.maxPerWindow {
		return false
	}
	s.counts[key]++
	return true
}
