package common

import (
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xrcuo/xrcuo-api/config"
)

func RequestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		endTime := time.Now()
		latency := endTime.Sub(startTime)

		if config.GetInstance().GetConfig().Log.RequestLog {
			logrus.WithFields(logrus.Fields{
				"method":     c.Request.Method,
				"path":       c.Request.URL.Path,
				"status":     c.Writer.Status(),
				"client_ip":  c.ClientIP(),
				"latency":    latency,
				"latency_ms": latency.Milliseconds(),
				"size":       c.Writer.Size(),
				"timestamp":  endTime.Format(time.RFC3339),
			}).Info("API请求")
		}
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		c.Writer.Header().Set("Access-Control-Max-Age", "3600")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, X-Response-Time")

		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval' https://cdn.jsdelivr.net; style-src 'self' 'unsafe-inline' https://cdn.jsdelivr.net https://fonts.googleapis.com; font-src 'self' data: https://fonts.gstatic.com; img-src * data:; connect-src 'self' https://cdn.jsdelivr.net")
		c.Writer.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Writer.Header().Set("Permissions-Policy", "geolocation=(self), camera=(), microphone=(), payment=()")
		c.Writer.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

type tokenBucket struct {
	capacity       float64
	rate           float64
	tokens         float64
	lastRefillTime time.Time
	mutex          sync.Mutex
}

func (tb *tokenBucket) take() bool {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefillTime).Seconds()

	newTokens := elapsed * tb.rate
	if newTokens > 0 {
		tb.tokens = math.Min(tb.tokens+newTokens, tb.capacity)
		tb.lastRefillTime = now
	}

	if tb.tokens >= 1 {
		tb.tokens--
		return true
	}

	return false
}

type rateLimiter struct {
	buckets         map[string]*tokenBucket
	mutex           sync.RWMutex
	capacity        float64
	rate            float64
	cleanupInterval time.Duration
	inactiveTimeout time.Duration
}

var globalRateLimiter *rateLimiter

func InitRateLimiter() {
	globalRateLimiter = &rateLimiter{
		buckets:         make(map[string]*tokenBucket),
		capacity:        config.GetRateLimitCapacity(),
		rate:            config.GetRateLimitRate(),
		cleanupInterval: 10 * time.Minute,
		inactiveTimeout: 30 * time.Minute,
	}

	logrus.Infof("速率限制器已初始化: 容量=%f, 速率=%f/秒",
		globalRateLimiter.capacity, globalRateLimiter.rate)

	go func() {
		ticker := time.NewTicker(globalRateLimiter.cleanupInterval)
		defer ticker.Stop()

		for range ticker.C {
			globalRateLimiter.cleanupInactiveBuckets()
		}
	}()
}

func (rl *rateLimiter) Allow(key string) bool {
	rl.mutex.RLock()
	tb, exists := rl.buckets[key]
	rl.mutex.RUnlock()

	if !exists {
		rl.mutex.Lock()
		if tb, exists = rl.buckets[key]; !exists {
			tb = &tokenBucket{
				capacity:       rl.capacity,
				rate:           rl.rate,
				tokens:         rl.capacity,
				lastRefillTime: time.Now(),
			}
			rl.buckets[key] = tb
		}
		rl.mutex.Unlock()
	}

	return tb.take()
}

func (rl *rateLimiter) cleanupInactiveBuckets() {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	for key, tb := range rl.buckets {
		tb.mutex.Lock()
		inactiveTime := now.Sub(tb.lastRefillTime)
		tb.mutex.Unlock()

		if inactiveTime > rl.inactiveTimeout {
			delete(rl.buckets, key)
		}
	}
}

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		if !globalRateLimiter.Allow(clientIP) {
			ErrorResponse(c, http.StatusTooManyRequests, 429, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}

		c.Next()
	}
}

func StatsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		method := c.Request.Method
		clientIP := c.ClientIP()

		c.Next()

		statusCode := c.Writer.Status()

		if GlobalStats != nil {
			go GlobalStats.RecordCall(path, method, clientIP, statusCode)
		}
	}
}

func PerformanceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		latency := time.Since(startTime)
		c.Writer.Header().Set("X-Response-Time", latency.String())
	}
}
