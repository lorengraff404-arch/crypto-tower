package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

// RateLimiter creates a rate limiting middleware
// Uses in-memory store (for production, use Redis)
func RateLimiter(rate string) gin.HandlerFunc {
	// Parse rate (e.g., "100-M" = 100 requests per minute)
	limiterRate, err := limiter.NewRateFromFormatted(rate)
	if err != nil {
		panic(err)
	}

	// Create memory store
	store := memory.NewStore()

	// Create limiter instance
	instance := limiter.New(store, limiterRate)

	// Return Gin middleware
	middleware := mgin.NewMiddleware(instance)
	
	return func(c *gin.Context) {
		middleware(c)
	}
}

// StrictRateLimiter for sensitive endpoints (auth, transactions)
func StrictRateLimiter() gin.HandlerFunc {
	return RateLimiter("10-M") // 10 requests per minute
}

// StandardRateLimiter for normal endpoints
func StandardRateLimiter() gin.HandlerFunc {
	return RateLimiter("100-M") // 100 requests per minute
}

// GenerousRateLimiter for read-heavy endpoints
func GenerousRateLimiter() gin.HandlerFunc {
	return RateLimiter("300-M") // 300 requests per minute
}

// CustomRateLimiter with custom configuration
func CustomRateLimiter(requests int64, period time.Duration) gin.HandlerFunc {
	rate := limiter.Rate{
		Period: period,
		Limit:  requests,
	}

	store := memory.NewStore()
	instance := limiter.New(store, rate)
	middleware := mgin.NewMiddleware(instance)

	return func(c *gin.Context) {
		middleware(c)
	}
}

// RateLimitByUserID limits based on user ID (requires auth)
func RateLimitByUserID(rate string) gin.HandlerFunc {
	limiterRate, err := limiter.NewRateFromFormatted(rate)
	if err != nil {
		panic(err)
	}

	store := memory.NewStore()
	instance := limiter.New(store, limiterRate, limiter.WithClientIPHeader("X-User-ID"))

	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		if userID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required",
			})
			c.Abort()
			return
		}

		// Set user ID as key for rate limiting
		c.Request.Header.Set("X-User-ID", string(rune(userID)))

		middleware := mgin.NewMiddleware(instance)
		middleware(c)
	}
}
