package ratelimit

import (
	"context"
	"testing"
	"time"
)

func TestWaitAllowsBurstThenThrottles(t *testing.T) {
	l := New(Config{RequestsPerMinute: 3})

	// Freeze the clock so timing is deterministic.
	base := time.Now()
	l.requests.now = func() time.Time { return base }
	l.requests.last = base

	// The full burst (3) must pass without blocking.
	for i := 0; i < 3; i++ {
		if err := l.Wait(context.Background(), 0); err != nil {
			t.Fatalf("burst call %d unexpectedly blocked: %v", i, err)
		}
	}

	// The next reservation must wait ~20s (3 per minute => 1 slot / 20s).
	d := l.requests.reserve(1)
	if d < 19*time.Second || d > 21*time.Second {
		t.Fatalf("expected ~20s wait after burst, got %v", d)
	}
}

func TestWaitRejectsRequestLargerThanTokenBudget(t *testing.T) {
	l := New(Config{TokensPerMinute: 10_000})
	if err := l.Wait(context.Background(), 10_001); err == nil {
		t.Fatal("expected an error for a request exceeding the per-minute token budget")
	}
}

func TestWaitHonoursContextCancellation(t *testing.T) {
	l := New(Config{RequestsPerMinute: 1})
	// Drain the single burst slot.
	if err := l.Wait(context.Background(), 0); err != nil {
		t.Fatalf("first call blocked: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // already cancelled; the second call would otherwise wait 60s
	if err := l.Wait(ctx, 0); err == nil {
		t.Fatal("expected context cancellation error while throttled")
	}
}

func TestEstimateTokens(t *testing.T) {
	if got := EstimateTokens(""); got != 1 {
		t.Fatalf("empty string should estimate to 1, got %d", got)
	}
	if got := EstimateTokens("12345678"); got != 2 { // 8 bytes / 4
		t.Fatalf("expected 2 tokens for 8 bytes, got %d", got)
	}
}

func TestNilDimensionsAreUnlimited(t *testing.T) {
	l := New(Config{}) // no limits at all
	for i := 0; i < 100; i++ {
		if err := l.Wait(context.Background(), 5_000); err != nil {
			t.Fatalf("unlimited limiter blocked on call %d: %v", i, err)
		}
	}
}
