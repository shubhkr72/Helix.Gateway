package ratelimiter

import (
	"context"
	"testing"
	"time"
)

func TestRedisLimiter(t *testing.T) {
	client, err := NewRedisClient("localhost:6379")
	if err != nil {
		t.Fatalf("failed to connect to redis: %v", err)
	}
	defer client.Close()

	cfg := Config{
		Capacity:    3,
		RefillRate:  1,
		KeyStrategy: "ip",
	}

	limiter := NewRedisLimiter(client, cfg)

	ctx := context.Background()
	key := "redis-test"

	client.Del(ctx, "ratelimit:"+key)

	// First 3 requests should be allowed.
	for i := 0; i < 3; i++ {
		result, err := limiter.Allow(ctx, key)
		if err != nil {
			t.Fatal(err)
		}

		if !result.Allowed {
			t.Fatalf("request %d should be allowed", i+1)
		}

		t.Logf(
			"Request %d -> Allowed=%v Remaining=%d RetryAfter=%v ResetAfter=%v",
			i+1,
			result.Allowed,
			result.Remaining,
			result.RetryAfter,
			result.ResetAfter,
		)
	}

	// Fourth request should be rejected.
	result, err := limiter.Allow(ctx, key)
	if err != nil {
		t.Fatal(err)
	}

	if result.Allowed {
		t.Fatal("expected rate limit exceeded")
	}

	t.Logf(
		"Rejected -> Remaining=%d RetryAfter=%v ResetAfter=%v",
		result.Remaining,
		result.RetryAfter,
		result.ResetAfter,
	)

	// Wait for refill.
	time.Sleep(2 * time.Second)

	result, err = limiter.Allow(ctx, key)
	if err != nil {
		t.Fatal(err)
	}

	if !result.Allowed {
		t.Fatal("expected token refill")
	}

	t.Logf(
		"After refill -> Allowed=%v Remaining=%d",
		result.Allowed,
		result.Remaining,
	)
}
