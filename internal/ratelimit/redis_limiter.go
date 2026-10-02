package ratelimit

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisLimiter struct {
	client *redis.Client
	script *redis.Script
}

func NewRedisLimiter(client *redis.Client) *RedisLimiter {
	return &RedisLimiter{
		client: client,
		script: redis.NewScript(tokenBucketScript),
	}
}

func (l *RedisLimiter) Allow(ctx context.Context, key string, cfg Config) (*Result, error) {
	now := float64(time.Now().UnixNano()) / 1e9 // Unix timestamp với phần thập phân

	res, err := l.script.Run(
		ctx,
		l.client,
		[]string{key},
		cfg.Capacity,
		cfg.RefillRate,
		now,
		1, // requested = 1 token
	).Result()

	if err != nil {
		return nil, err
	}

	// Parse kết quả {allowed, remaining, retry_after}
	arr, ok := res.([]interface{})
	if !ok || len(arr) != 3 {
		// Fallback: cho phép nếu không parse được (fail-open)
		return &Result{Allowed: true, Remaining: cfg.Capacity, Limit: cfg.Capacity}, nil
	}

	allowed := toInt64(arr[0]) == 1
	remaining := int(toInt64(arr[1]))
	retryAfter := time.Duration(toInt64(arr[2])) * time.Second

	return &Result{
		Allowed:    allowed,
		Remaining:  remaining,
		RetryAfter: retryAfter,
		Limit:      cfg.Capacity,
	}, nil
}

func toInt64(v interface{}) int64 {
	switch x := v.(type) {
	case int64:
		return x
	case int:
		return int64(x)
	case float64:
		return int64(x)
	}
	return 0
}
