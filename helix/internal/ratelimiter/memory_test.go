package ratelimiter

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestMemoryLimiterCapacity(t *testing.T) {
	clock := NewFakeClock(time.Unix(0, 0))

	limiter := NewMemoryLimiter(
		Config{
			Capacity:   5,
			RefillRate: 1,
		},
		clock,
	)

	for i := 0; i < 5; i++ {
		allowed, err := limiter.Allow(context.Background(), "user1")
		if err != nil {
			t.Fatal(err)
		}

		if !allowed {
			t.Fatalf("request %d should be allowed", i+1)
		}
	}

	allowed, err := limiter.Allow(context.Background(), "user1")
	if err != nil {
		t.Fatal(err)
	}

	if allowed {
		t.Fatal("expected rate limit exceeded")
	}
}

func TestMemoryLimiterRefill(t *testing.T) {
	clock := NewFakeClock(time.Unix(0, 0))

	limiter := NewMemoryLimiter(
		Config{
			Capacity:   5,
			RefillRate: 2,
		},
		clock,
	)

	for i := 0; i < 5; i++ {
		limiter.Allow(context.Background(), "user1")
	}

	allowed, _ := limiter.Allow(context.Background(), "user1")
	if allowed {
		t.Fatal("expected limit exceeded")
	}

	clock.Advance(time.Second)

	allowed, _ = limiter.Allow(context.Background(), "user1")
	if !allowed {
		t.Fatal("expected one token after refill")
	}

	allowed, _ = limiter.Allow(context.Background(), "user1")
	if !allowed {
		t.Fatal("expected second token after refill")
	}

	allowed, _ = limiter.Allow(context.Background(), "user1")
	if allowed {
		t.Fatal("expected bucket empty again")
	}
}

func TestMemoryLimiterSeparateBuckets(t *testing.T) {
	clock := NewFakeClock(time.Unix(0, 0))

	limiter := NewMemoryLimiter(
		Config{
			Capacity:   2,
			RefillRate: 1,
		},
		clock,
	)

	limiter.Allow(context.Background(), "user1")
	limiter.Allow(context.Background(), "user1")

	allowed, _ := limiter.Allow(context.Background(), "user1")
	if allowed {
		t.Fatal("user1 should be limited")
	}

	allowed, _ = limiter.Allow(context.Background(), "user2")
	if !allowed {
		t.Fatal("user2 should have independent bucket")
	}
}

func TestMemoryLimiterCapacityLimit(t *testing.T) {
	clock := NewFakeClock(time.Unix(0, 0))

	limiter := NewMemoryLimiter(
		Config{
			Capacity:   5,
			RefillRate: 10,
		},
		clock,
	)

	limiter.Allow(context.Background(), "user1")

	clock.Advance(10 * time.Second)

	for i := 0; i < 5; i++ {
		allowed, _ := limiter.Allow(context.Background(), "user1")
		if !allowed {
			t.Fatal("expected token available")
		}
	}

	allowed, _ := limiter.Allow(context.Background(), "user1")
	if allowed {
		t.Fatal("bucket should never exceed capacity")
	}
}

func TestMemoryLimiterConcurrent(t *testing.T) {
	clock := NewFakeClock(time.Unix(0, 0))

	limiter := NewMemoryLimiter(
		Config{
			Capacity:   100,
			RefillRate: 0,
		},
		clock,
	)

	var wg sync.WaitGroup

	success := 0
	var mu sync.Mutex

	for i := 0; i < 500; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			ok, _ := limiter.Allow(context.Background(), "user")

			if ok {
				mu.Lock()
				success++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	if success != 100 {
		t.Fatalf("expected 100 successful requests, got %d", success)
	}
}
