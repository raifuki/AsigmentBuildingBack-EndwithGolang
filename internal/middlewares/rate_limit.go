package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/raifuki/task-management/internal/ratelimit"
	"github.com/raifuki/task-management/pkg/response"
)

// RateLimitByIP giới hạn theo IP (cho public endpoints)
func RateLimitByIP(limiter ratelimit.Limiter, cfg ratelimit.Config, prefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := fmt.Sprintf("ratelimit:%s:ip:%s", prefix, c.ClientIP())
		applyRateLimit(c, limiter, key, cfg)
	}
}

// RateLimitByUser giới hạn theo UserID (cho private endpoints)
// Phải đặt SAU middleware Auth()
func RateLimitByUser(limiter ratelimit.Limiter, cfg ratelimit.Config, prefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get(CtxUserID)
		if !exists {
			// Chưa auth → fallback theo IP
			key := fmt.Sprintf("ratelimit:%s:ip:%s", prefix, c.ClientIP())
			applyRateLimit(c, limiter, key, cfg)
			return
		}
		key := fmt.Sprintf("ratelimit:%s:user:%v", prefix, userID)
		applyRateLimit(c, limiter, key, cfg)
	}
}

// applyRateLimit là logic chung
func applyRateLimit(c *gin.Context, limiter ratelimit.Limiter, key string, cfg ratelimit.Config) {
	res, err := limiter.Allow(c.Request.Context(), key, cfg)
	if err != nil {
		// Lỗi Redis → fail-open (cho qua) để không chặn nhầm user
		log.Printf("⚠️ Rate limiter error: %v", err)
		c.Next()
		return
	}

	// Set headers chuẩn (giống GitHub API)
	c.Header("X-RateLimit-Limit", strconv.Itoa(res.Limit))
	c.Header("X-RateLimit-Remaining", strconv.Itoa(res.Remaining))

	if !res.Allowed {
		c.Header("Retry-After", strconv.Itoa(int(res.RetryAfter.Seconds())))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(res.RetryAfter).Unix(), 10))

		response.Error(c, http.StatusTooManyRequests, "Too many requests, please slow down", gin.H{
			"retry_after_seconds": int(res.RetryAfter.Seconds()),
		})
		c.Abort()
		return
	}

	c.Next()
}
