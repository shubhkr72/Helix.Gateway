package ratelimiter

import (
	"context"
	"time"
)

type Limiter interface {
	Allow(ctx context.Context, key string) (Result, error)
}

type Result struct {
	Allowed    bool
	Limit      int64
	Remaining  int64
	RetryAfter time.Duration
	ResetAfter time.Duration
}
