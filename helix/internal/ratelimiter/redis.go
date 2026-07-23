package ratelimiter

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisLimiter struct {
	client   *redis.Client
	capacity float64
	refill   float64
	script   *redis.Script
}

func NewRedisClient(addr string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return client, nil
}

func NewRedisLimiter(client *redis.Client, cfg Config) *RedisLimiter {
	return &RedisLimiter{
		client:   client,
		capacity: cfg.Capacity,
		refill:   cfg.RefillRate,
		script:   redis.NewScript(luaTokenBucket),
	}
}

func (r *RedisLimiter) Allow(ctx context.Context, key string) (Result, error) {
	if r.client == nil {
		return Result{}, errors.New("redis unavailable")
	}

	result, err := r.script.Run(
		ctx,
		r.client,
		[]string{fmt.Sprintf("ratelimit:%s", key)},
		r.capacity,
		r.refill,
	).Result()

	if err != nil {
		return Result{}, err
	}

	values, ok := result.([]interface{})
	if !ok {
		return Result{}, errors.New("unexpected redis response type")
	}

	if len(values) != 4 {
		return Result{}, errors.New("unexpected redis response length")
	}

	allowed, ok := values[0].(int64)
	if !ok {
		return Result{}, errors.New("invalid allowed value")
	}

	remaining, ok := values[1].(int64)
	if !ok {
		return Result{}, errors.New("invalid remaining value")
	}

	retryAfter, ok := values[2].(int64)
	if !ok {
		return Result{}, errors.New("invalid retry_after value")
	}

	resetAfter, ok := values[3].(int64)
	if !ok {
		return Result{}, errors.New("invalid reset_after value")
	}

	return Result{
		Allowed:    allowed == 1,
		Limit:      int64(r.capacity),
		Remaining:  remaining,
		RetryAfter: time.Duration(retryAfter) * time.Second,
		ResetAfter: time.Duration(resetAfter) * time.Second,
	}, nil
}
