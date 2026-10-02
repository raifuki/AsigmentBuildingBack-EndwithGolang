package ratelimit

import (
	"context"
	"sync"
	"time"
)

type memBucket struct {
	tokens     float64
	lastRefill time.Time
}

type MemoryLimiter struct {
	mu      sync.Mutex
	buckets map[string]*memBucket
}

func NewMemoryLimiter() *MemoryLimiter {
	m := &MemoryLimiter{buckets: make(map[string]*memBucket)}
	go m.cleanup()
	return m
}

func (m *MemoryLimiter) Allow(ctx context.Context, key string, cfg Config) (*Result, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	b, ok := m.buckets[key]
	if !ok {
		b = &memBucket{tokens: float64(cfg.Capacity), lastRefill: now}
		m.buckets[key] = b
	}

	elapsed := now.Sub(b.lastRefill).Seconds()
	b.tokens = minFloat64(float64(cfg.Capacity), b.tokens+elapsed*cfg.RefillRate)
	b.lastRefill = now

	if b.tokens >= 1 {
		b.tokens--
		return &Result{Allowed: true, Remaining: int(b.tokens), Limit: cfg.Capacity}, nil
	}

	retryAfter := time.Duration((1-b.tokens)/cfg.RefillRate) * time.Second
	return &Result{
		Allowed:    false,
		Remaining:  0,
		RetryAfter: retryAfter,
		Limit:      cfg.Capacity,
	}, nil
}

func (m *MemoryLimiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		m.mu.Lock()
		for k, b := range m.buckets {
			if time.Since(b.lastRefill) > 10*time.Minute {
				delete(m.buckets, k)
			}
		}
		m.mu.Unlock()
	}
}

func minFloat64(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// FallbackLimiter tries primary limiter first, and falls back to memory limiter when primary fails.
type FallbackLimiter struct {
	primary  Limiter
	fallback Limiter
}

func NewFallbackLimiter(primary Limiter, fallback Limiter) *FallbackLimiter {
	return &FallbackLimiter{primary: primary, fallback: fallback}
}

func (f *FallbackLimiter) Allow(ctx context.Context, key string, cfg Config) (*Result, error) {
	res, err := f.primary.Allow(ctx, key, cfg)
	if err != nil {
		return f.fallback.Allow(ctx, key, cfg)
	}
	return res, nil
}
