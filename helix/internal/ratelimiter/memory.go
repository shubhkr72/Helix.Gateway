package ratelimiter

import (
	"context"
	"sync"
)

type MemoryLimiter struct {
	mu       sync.Mutex
	buckets  map[string]*bucket
	capacity float64
	refill   float64
	clock    Clock
}

func NewMemoryLimiter(cfg Config, clock Clock) *MemoryLimiter {
	if clock == nil {
		clock = RealClock{}
	}

	return &MemoryLimiter{
		buckets:  make(map[string]*bucket),
		capacity: cfg.Capacity,
		refill:   cfg.RefillRate,
		clock:    clock,
	}
}

func (m *MemoryLimiter) Allow(ctx context.Context, key string) (bool, error) {
	_ = ctx

	m.mu.Lock()
	defer m.mu.Unlock()

	now := m.clock.Now()

	b, ok := m.buckets[key]
	if !ok {
		b = &bucket{
			tokens: m.capacity,
			last:   now,
		}
		m.buckets[key] = b
	}

	elapsed := now.Sub(b.last).Seconds()

	if elapsed > 0 {
		b.tokens += elapsed * m.refill

		if b.tokens > m.capacity {
			b.tokens = m.capacity
		}

		b.last = now
	}

	if b.tokens < 1 {
		return false, nil
	}

	b.tokens--

	return true, nil
}
