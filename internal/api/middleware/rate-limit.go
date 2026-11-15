package middleware

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"justcallmesu.com/golang-base/internal/common/apperror"
)

// RateLimit allows `limit` requests per `period` for each client IP,
// e.g. RateLimit(10, time.Minute). Each call has its own counter.
// Call it while registering routes (once at startup), never inside a handler.
func RateLimit(limit int64, period time.Duration) gin.HandlerFunc {
	rateLimiter := limiter.New(memory.NewStore(), limiter.Rate{Limit: limit, Period: period})

	return func(c *gin.Context) {
		result, err := rateLimiter.Get(c.Request.Context(), c.ClientIP())

		if err != nil {
			_ = c.Error(apperror.Internal(fmt.Errorf("rate limiter: %w", err)))
			c.Abort()
			return
		}

		c.Header("X-RateLimit-Limit", strconv.FormatInt(result.Limit, 10))
		c.Header("X-RateLimit-Remaining", strconv.FormatInt(result.Remaining, 10))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(result.Reset, 10))

		if result.Reached {
			secondsUntilReset := math.Ceil(time.Until(time.Unix(result.Reset, 0)).Seconds())
			c.Header("Retry-After", strconv.Itoa(max(1, int(secondsUntilReset))))

			_ = c.Error(apperror.TooManyRequests("too many requests, please try again later"))
			c.Abort()
			return
		}

		c.Next()

	}
}
