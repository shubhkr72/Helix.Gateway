package ratelimiter

import (
	"context"
	"errors"
	"fmt"

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

func (r *RedisLimiter) Allow(ctx context.Context, key string) (bool, error) {
	result, err := r.script.Run(
		ctx,
		r.client,
		[]string{fmt.Sprintf("ratelimit:%s", key)},
		r.capacity,
		r.refill,
	).Result()

	if err != nil {
		return false, err
	}

	value, ok := result.(int64)
	if !ok {
		return false, errors.New("unexpected redis response")
	}

	return value == 1, nil
}
