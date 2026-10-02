package ratelimit

import (
	"context"
	"time"
)

// Result trả về từ limiter
type Result struct {
	Allowed    bool          // cho phép request hay không
	Remaining  int           // số token còn lại
	RetryAfter time.Duration // thời gian cần chờ (nếu bị chặn)
	Limit      int           // capacity (để set header X-RateLimit-Limit)
}

// Config cho một bucket
type Config struct {
	Capacity   int           // số token tối đa
	RefillRate float64       // token/giây
	Window     time.Duration // không dùng, chỉ để tham khảo
}

// Limiter interface — dễ swap giữa Redis / in-memory
type Limiter interface {
	Allow(ctx context.Context, key string, cfg Config) (*Result, error)
}

// Preset config cho các loại endpoint khác nhau
var (
	// Auth endpoints: 5 request/phút (chống brute-force)
	AuthLimit = Config{
		Capacity:   5,
		RefillRate: 5.0 / 60.0, // 5 per minute
	}

	// Global API: 60 req/phút
	GlobalLimit = Config{
		Capacity:   60,
		RefillRate: 1.0, // 1 per second
	}

	// Write endpoints (POST/PUT/DELETE): 20 req/phút
	WriteLimit = Config{
		Capacity:   20,
		RefillRate: 20.0 / 60.0,
	}

	// Read endpoints (GET): 100 req/phút
	ReadLimit = Config{
		Capacity:   100,
		RefillRate: 100.0 / 60.0,
	}
)
