package agent

import (
	"testing"
	"time"
)

// --- pruneOld ---

func TestPruneOld(t *testing.T) {
	now := time.Now()
	old := now.Add(-2 * time.Minute)
	recent := now.Add(-30 * time.Second)
	cutoff := now.Add(-time.Minute)

	tests := []struct {
		name    string
		times   []time.Time
		wantLen int
	}{
		{"empty slice", nil, 0},
		{"all old", []time.Time{old, old}, 0},
		{"all recent", []time.Time{recent, recent}, 2},
		{"mixed", []time.Time{old, old, recent, recent}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pruneOld(tt.times, cutoff)
			if len(got) != tt.wantLen {
				t.Errorf("pruneOld returned %d items, want %d", len(got), tt.wantLen)
			}
		})
	}
}

// --- newRateLimiter ---

func TestNewRateLimiter(t *testing.T) {
	rl := newRateLimiter(10, 5)
	if rl.maxToolCallsPerMinute != 10 {
		t.Errorf("maxToolCallsPerMinute = %d, want 10", rl.maxToolCallsPerMinute)
	}
	if rl.maxRequestsPerMinute != 5 {
		t.Errorf("maxRequestsPerMinute = %d, want 5", rl.maxRequestsPerMinute)
	}
	if len(rl.toolCallTimes) != 0 {
		t.Error("toolCallTimes should be empty")
	}
	if len(rl.requestTimes) != 0 {
		t.Error("requestTimes should be empty")
	}
}

// --- checkToolCall ---

func TestCheckToolCall_Unlimited(t *testing.T) {
	rl := newRateLimiter(0, 0)
	for i := 0; i < 100; i++ {
		if err := rl.checkToolCall(); err != nil {
			t.Fatalf("checkToolCall should never fail with limit=0, got %v", err)
		}
	}
}

func TestCheckToolCall_WithinLimit(t *testing.T) {
	rl := newRateLimiter(5, 0)
	for i := 0; i < 5; i++ {
		if err := rl.checkToolCall(); err != nil {
			t.Fatalf("call %d should succeed: %v", i, err)
		}
	}
}

func TestCheckToolCall_ExceedsLimit(t *testing.T) {
	rl := newRateLimiter(3, 0)
	for i := 0; i < 3; i++ {
		if err := rl.checkToolCall(); err != nil {
			t.Fatalf("call %d should succeed: %v", i, err)
		}
	}
	if err := rl.checkToolCall(); err == nil {
		t.Error("4th call should exceed limit")
	}
}

// --- checkRequest ---

func TestCheckRequest_Unlimited(t *testing.T) {
	rl := newRateLimiter(0, 0)
	for i := 0; i < 100; i++ {
		if err := rl.checkRequest(); err != nil {
			t.Fatalf("checkRequest should never fail with limit=0, got %v", err)
		}
	}
}

func TestCheckRequest_WithinLimit(t *testing.T) {
	rl := newRateLimiter(0, 5)
	for i := 0; i < 5; i++ {
		if err := rl.checkRequest(); err != nil {
			t.Fatalf("request %d should succeed: %v", i, err)
		}
	}
}

func TestCheckRequest_ExceedsLimit(t *testing.T) {
	rl := newRateLimiter(0, 2)
	for i := 0; i < 2; i++ {
		if err := rl.checkRequest(); err != nil {
			t.Fatalf("request %d should succeed: %v", i, err)
		}
	}
	if err := rl.checkRequest(); err == nil {
		t.Error("3rd request should exceed limit")
	}
}
