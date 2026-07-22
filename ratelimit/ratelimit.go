package ratelimit

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Config struct {
	RequestsPerMinute int
	TokensPerMinute   int
}

type Limiter struct {
	requests *bucket // nil => no request limit
	tokens   *bucket // nil => no token limit
	maxTok   int     // token burst, used to reject impossible requests early
}

func New(cfg Config) *Limiter {
	l := &Limiter{}
	if cfg.RequestsPerMinute > 0 {
		l.requests = newBucket(float64(cfg.RequestsPerMinute)/60.0, cfg.RequestsPerMinute)
	}
	if cfg.TokensPerMinute > 0 {
		l.tokens = newBucket(float64(cfg.TokensPerMinute)/60.0, cfg.TokensPerMinute)
		l.maxTok = cfg.TokensPerMinute
	}
	return l
}

func (l *Limiter) Wait(ctx context.Context, tokens int) error {
	if l.tokens != nil && tokens > l.maxTok {
		return fmt.Errorf("ratelimit: request of %d tokens exceeds the %d tokens/min budget; split the input", tokens, l.maxTok)
	}
	if l.requests != nil {
		if err := l.requests.wait(ctx, 1); err != nil {
			return err
		}
	}
	if l.tokens != nil && tokens > 0 {
		if err := l.tokens.wait(ctx, tokens); err != nil {
			return err
		}
	}
	return nil
}

// EstimateTokens returns a rough token count for text using a ~4-bytes-per-token
// heuristic. It is intentionally cheap and slightly conservative for budgeting;
// it is not a real tokenizer. Keep a 429-aware retry as the correctness backstop
// for when the estimate runs low.
func EstimateTokens(text string) int {
	n := (len(text) + 3) / 4
	if n < 1 {
		return 1
	}
	return n
}

// bucket is a classic token bucket. Reservations may drive the balance negative;
// the deficit is exactly the amount of future refill a caller must wait for.
type bucket struct {
	mu       sync.Mutex
	capacity float64
	tokens   float64
	rate     float64 // tokens per second
	last     time.Time
	now      func() time.Time // injectable for tests
}

func newBucket(ratePerSec float64, burst int) *bucket {
	return &bucket{
		capacity: float64(burst),
		tokens:   float64(burst),
		rate:     ratePerSec,
		last:     time.Now(),
		now:      time.Now,
	}
}

// reserve refills according to elapsed time, deducts n, and returns how long the
// caller must wait before those n tokens are actually available (0 if already).
func (b *bucket) reserve(n int) time.Duration {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := b.now()
	b.tokens += now.Sub(b.last).Seconds() * b.rate
	if b.tokens > b.capacity {
		b.tokens = b.capacity
	}
	b.last = now

	b.tokens -= float64(n)
	if b.tokens >= 0 {
		return 0
	}
	return time.Duration(-b.tokens / b.rate * float64(time.Second))
}

func (b *bucket) wait(ctx context.Context, n int) error {
	delay := b.reserve(n)
	if delay <= 0 {
		return nil
	}
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
