package ratelimiter

import (
	"context"
	"testing"
	"time"
)

func TestRedisLimiter(t *testing.T) {
	cfg := Config{
		Capacity:    3,
		RefillRate:  1,
		KeyStrategy: "ip",
	}

	limiter, err := NewRedisLimiter("localhost:6379", cfg)
	if err != nil {
		t.Fatalf("failed to connect to redis: %v", err)
	}
	defer limiter.Close()

	ctx := context.Background()
	key := "redis-test"

	limiter.client.Del(ctx, "ratelimit:"+key)

	for i := 0; i < 3; i++ {
		ok, err := limiter.Allow(ctx, key)
		if err != nil {
			t.Fatal(err)
		}

		if !ok {
			t.Fatalf("request %d should be allowed", i+1)
		}
	}

	ok, err := limiter.Allow(ctx, key)
	if err != nil {
		t.Fatal(err)
	}

	if ok {
		t.Fatal("expected rate limit exceeded")
	}

	time.Sleep(2 * time.Second)

	ok, err = limiter.Allow(ctx, key)
	if err != nil {
		t.Fatal(err)
	}

	if !ok {
		t.Fatal("expected token refill")
	}
}
