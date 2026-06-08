package middleware

import (
	"ecommerce/pkg/logger"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
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

		// log.Printf("[API] %s | %d | %s %s | %v", reqID, c.Writer.Status(), c.Request.Method, c.Request.URL.Path, latency)
		logger.Log.Info("Incoming Request",
			zap.Any("request_id", reqID),
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Duration("latency", latency),
			zap.String("ip", c.ClientIP()),
		)
	}
}

const tokenBucketScript = `
local tokens_key = KEYS[1]
local timestamp_key = KEYS[2]

local rate = tonumber(ARGV[1])
local capacity = tonumber(ARGV[2])
local now = tonumber(ARGV[3])
local requested = 1

local fill_time = capacity / rate
local ttl = math.floor(fill_time * 2)

-- Ambil sisa token saat ini
local last_tokens = tonumber(redis.call("GET", tokens_key))
if last_tokens == nil then
	last_tokens = capacity
end

-- Ambil waktu terakhir kali token diisi
local last_refreshed = tonumber(redis.call("GET", timestamp_key))
if last_refreshed == nil then
	last_refreshed = 0
end

-- Hitung berapa banyak token baru yang terisi sejak request terakhir
local delta = math.max(0, now - last_refreshed)
local filled_tokens = math.min(capacity, last_tokens + (delta * rate))
local allowed = filled_tokens >= requested

if allowed then
	local new_tokens = filled_tokens - requested
	redis.call("SETEX", tokens_key, ttl, new_tokens)
	redis.call("SETEX", timestamp_key, ttl, now)
	return 1
else
	redis.call("SETEX", tokens_key, ttl, filled_tokens)
	redis.call("SETEX", timestamp_key, ttl, now)
	return 0
end
`

func RateLimiterTokenBucket(redisClient *redis.Client, ratePerSec float64, capacity int) gin.HandlerFunc {
	script := redis.NewScript(tokenBucketScript)

	return func(c *gin.Context) {
		ctx := c.Request.Context()
		ip := c.ClientIP()

		tokensKey := fmt.Sprintf("rate_limit:tokens:%s", ip)
		timestampKey := fmt.Sprintf("rate_limit:ts:%s", ip)
		now := time.Now().Unix()

		result, err := script.Run(ctx, redisClient, []string{tokensKey, timestampKey}, ratePerSec, capacity, now).Result()

		if err != nil {
			// log.Printf("⚠️ Redis error saat rate limiting: %v", err)
			logger.Log.Warn("Redis error saat rate limiting", zap.Error(err))
			c.Next()
			return
		}

		allowed := result.(int64) == 1

		if !allowed {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Terlalu banyak permintaan. Silakan tunggu beberapa saat."})
			c.Abort()
			return
		}

		c.Next()
	}
}
