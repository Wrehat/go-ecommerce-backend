package middleware

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/time/rate"
)

// Todo : middleware untuk request id
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.GetHeader("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}

		c.Set("request_id", reqID)
		c.Header("X-Request-ID", reqID)
		c.Next()
	}
}

// Todo : middleware untuk logger
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)
		reqID, _ := c.Get("request_id")

		log.Printf("[API] %s | %d | %s %s | %v", reqID, c.Writer.Status(), c.Request.Method, c.Request.URL.Path, latency)
	}
}

// Todo : Rate limiter middleware
var (
	mu       sync.Mutex
	limiters = make(map[string]*rate.Limiter)
)

func getLimiter(ip string, r rate.Limit, b int) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(r, b)
		limiters[ip] = limiter
	}
	return limiter
}

func RateLimiter(r rate.Limit, b int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := getLimiter(ip, r, b)

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too Many Requests. Silakan tunggu."})
			c.Abort()
			return
		}
		c.Next()
	}
}
