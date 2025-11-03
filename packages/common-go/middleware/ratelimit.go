package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/streamverse/common-go/errors"
)

// RateLimiter implements rate limiting using Redis
type RateLimiter struct {
	client *redis.Client
	limit  int
	window time.Duration
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(client *redis.Client, limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		client: client,
		limit:  limit,
		window: window,
	}
}

// RateLimit middleware
func (rl *RateLimiter) RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get client identifier (IP or user ID)
		identifier := c.ClientIP()
		if userID, exists := c.Get("user_id"); exists {
			identifier = fmt.Sprintf("user:%s", userID)
		}

		key := fmt.Sprintf("rate_limit:%s", identifier)
		ctx := context.Background()

		// Get current count
		count, err := rl.client.Incr(ctx, key).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, errors.NewInternalError("Rate limit check failed"))
			c.Abort()
			return
		}

		// Set expiration on first request
		if count == 1 {
			rl.client.Expire(ctx, key, rl.window)
		}

		// Check if limit exceeded
		if count > int64(rl.limit) {
			c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", rl.limit))
			c.Header("X-RateLimit-Remaining", "0")
			c.JSON(http.StatusTooManyRequests, errors.NewAppError(
				errors.ErrorCodeServiceUnavailable,
				"Rate limit exceeded",
				http.StatusTooManyRequests,
			))
			c.Abort()
			return
		}

		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", rl.limit))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", rl.limit-int(count)))

		c.Next()
	}
}

