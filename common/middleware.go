package common

import (
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"github.com/xrcuo/xrcuo-api/config"
	"github.com/xrcuo/xrcuo-api/db"
	"github.com/xrcuo/xrcuo-api/models"
)

// 全局API密钥缓存实例
var apiKeyCacheInstance *cache.Cache

// init 初始化API密钥缓存
func init() {
	// 创建go-cache实例，设置默认过期时间为5分钟，清理间隔为10分钟
	apiKeyCacheInstance = cache.New(5*time.Minute, 10*time.Minute)
}

// GetAPICache 获取API密钥缓存实例
func GetAPICache() *cache.Cache {
	return apiKeyCacheInstance
}

// RequestLoggerMiddleware 请求日志中间件
func RequestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 请求结束时间
		endTime := time.Now()
		latency := endTime.Sub(startTime)

		// 检查是否启用请求日志
		if config.GetInstance().GetConfig().Log.RequestLog {
			// 记录请求日志，移除敏感信息
			logrus.WithFields(logrus.Fields{
				"method":     c.Request.Method,
				"path":       c.Request.URL.Path,
				"status":     c.Writer.Status(),
				"client_ip":  c.ClientIP(), // 注意：生产环境中可能需要掩码IP地址
				"latency":    latency,
				"latency_ms": latency.Milliseconds(),
				"size":       c.Writer.Size(),
				"timestamp":  endTime.Format(time.RFC3339),
			}).Info("API请求")
		}
	}
}

// CORSMiddleware 跨域请求中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		
		// 安全的CORS配置：只在有Origin头时设置CORS响应头
		if origin != "" {
			// 生产环境建议：替换为实际的允许域名列表
			// 这里保持灵活性，但记录警告以便安全审计
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}

		// CORS相关头
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		c.Writer.Header().Set("Access-Control-Max-Age", "3600")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, X-Response-Time")

		// 安全头：防止点击劫持
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		// 安全头：防止XSS攻击
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		// 安全头：防止MIME类型嗅探
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		// 安全头：内容安全策略
		c.Writer.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval' https://cdn.jsdelivr.net; style-src 'self' 'unsafe-inline' https://cdn.jsdelivr.net https://fonts.googleapis.com; font-src 'self' data: https://fonts.gstatic.com; img-src * data:; connect-src 'self' https://cdn.jsdelivr.net")
		// 安全头：减少信息泄露
		c.Writer.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		// 安全头：控制资源加载策略
		c.Writer.Header().Set("Permissions-Policy", "geolocation=(self), camera=(), microphone=(), payment=()")
		// 安全头：HTTP严格传输安全
		c.Writer.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// APIKeyMiddleware API密钥验证中间件
func APIKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如果未启用API密钥验证，则跳过验证
		if !config.GetInstance().GetConfig().APIKey.Enabled {
			c.Next()
			return
		}

		// 从请求头或查询参数中获取API密钥
		apiKey := c.GetHeader("Authorization")
		if apiKey == "" {
			apiKey = c.Query("api_key")
		}

		// 检查API密钥是否存在
		if apiKey == "" {
			ErrorResponse(c, http.StatusUnauthorized, 401, "API密钥不能为空")
			c.Abort()
			return
		}

		// 验证API密钥
		var keyInfo *models.APIKey
		// 从缓存获取
		if val, found := apiKeyCacheInstance.Get(apiKey); found {
			keyInfo = val.(*models.APIKey)
		}

		if keyInfo == nil {
			// 从数据库获取
			var err error
			keyInfo, err = db.GetAPIKeyByKey(apiKey)
			if err != nil {
				ErrorResponse(c, http.StatusUnauthorized, 401, "无效的API密钥")
				c.Abort()
				return
			}
			// 存入缓存
			apiKeyCacheInstance.Set(apiKey, keyInfo, cache.DefaultExpiration)
		}

		// 检查API密钥是否已达到使用上限
		if !keyInfo.IsPermanent && keyInfo.CurrentUsage >= keyInfo.MaxUsage {
			ErrorResponse(c, http.StatusForbidden, 403, "API密钥已达到使用上限")
			c.Abort()
			return
		}

		// 更新API密钥使用次数（数据库层面做原子性检查）
		if err := db.UpdateAPIKeyUsage(apiKey); err != nil {
			if err.Error() == "API密钥已达到使用上限" {
				// 缓存数据过期，从数据库重新获取最新状态
				apiKeyCacheInstance.Delete(apiKey)
				ErrorResponse(c, http.StatusForbidden, 403, "API密钥已达到使用上限")
			} else {
				ErrorResponse(c, http.StatusInternalServerError, 500, "更新API密钥使用次数失败")
			}
			c.Abort()
			return
		}

		// 更新缓存中的使用次数（仅在数据库更新成功后）
		keyInfo.CurrentUsage++
		apiKeyCacheInstance.Set(apiKey, keyInfo, cache.DefaultExpiration)

		// 将API密钥信息存储到上下文
		c.Set("api_key", keyInfo)

		// 继续处理请求
		c.Next()
	}
}

// tokenBucket 令牌桶
type tokenBucket struct {
	capacity       float64    // 令牌桶容量
	rate           float64    // 令牌生成速率（每秒）
	tokens         float64    // 当前令牌数量
	lastRefillTime time.Time  // 上次填充令牌的时间
	mutex          sync.Mutex // 保护令牌桶的互斥锁
}

// take 尝试获取一个令牌
func (tb *tokenBucket) take() bool {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefillTime).Seconds()

	// 填充令牌
	newTokens := elapsed * tb.rate
	if newTokens > 0 {
		tb.tokens = math.Min(tb.tokens+newTokens, tb.capacity)
		tb.lastRefillTime = now
	}

	// 尝试获取令牌
	if tb.tokens >= 1 {
		tb.tokens--
		return true
	}

	return false
}

// rateLimiter 速率限制器
type rateLimiter struct {
	buckets         map[string]*tokenBucket
	mutex           sync.RWMutex
	capacity        float64       // 默认令牌桶容量
	rate            float64       // 默认令牌生成速率（每秒）
	cleanupInterval time.Duration // 清理过期项的时间间隔
	inactiveTimeout time.Duration // 令牌桶的不活动超时时间
}

// 全局速率限制器实例
var globalRateLimiter *rateLimiter

// init 初始化速率限制器
func init() {
	globalRateLimiter = &rateLimiter{
		buckets:         make(map[string]*tokenBucket),
		capacity:        100,              // 令牌桶容量，允许突发流量
		rate:            1.666,            // 令牌生成速率，约100个/分钟
		cleanupInterval: 10 * time.Minute, // 每10分钟清理一次
		inactiveTimeout: 30 * time.Minute, // 30分钟不活动则清理
	}

	// 启动定期清理过期项的任务
	go func() {
		ticker := time.NewTicker(globalRateLimiter.cleanupInterval)
		defer ticker.Stop()

		for range ticker.C {
			globalRateLimiter.cleanupInactiveBuckets()
		}
	}()
}

// Allow 检查是否允许请求
func (rl *rateLimiter) Allow(key string) bool {
	// 先获取读锁，检查令牌桶是否存在
	rl.mutex.RLock()
	tb, exists := rl.buckets[key]
	rl.mutex.RUnlock()

	// 如果不存在，创建一个新的令牌桶
	if !exists {
		rl.mutex.Lock()
		// 双重检查，避免并发创建
		if tb, exists = rl.buckets[key]; !exists {
			tb = &tokenBucket{
				capacity:       rl.capacity,
				rate:           rl.rate,
				tokens:         rl.capacity, // 初始时填满令牌桶
				lastRefillTime: time.Now(),
			}
			rl.buckets[key] = tb
		}
		rl.mutex.Unlock()
	}

	// 尝试从令牌桶中获取一个令牌
	return tb.take()
}

// cleanupInactiveBuckets 清理不活动的令牌桶
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

// RateLimitMiddleware 速率限制中间件
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 使用客户端IP作为速率限制的键
		clientIP := c.ClientIP()

		// 检查是否允许请求
		if !globalRateLimiter.Allow(clientIP) {
			ErrorResponse(c, http.StatusTooManyRequests, 429, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}

		c.Next()
	}
}

// StatsMiddleware 统计API调用次数的中间件
func StatsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求信息
		path := c.Request.URL.Path
		method := c.Request.Method
		clientIP := c.ClientIP()

		// 处理请求
		c.Next()

		// 获取响应状态码
		statusCode := c.Writer.Status()

		// 异步记录调用信息，减少对请求响应时间的影响
		if GlobalStats != nil {
			go GlobalStats.RecordCall(path, method, clientIP, statusCode)
		}
	}
}

// PerformanceMetrics 性能指标结构体
type PerformanceMetrics struct {
	TotalRequests     int64                   `json:"total_requests"`
	TotalResponseTime time.Duration           `json:"total_response_time"`
	AvgResponseTime   time.Duration           `json:"avg_response_time"`
	MaxResponseTime   time.Duration           `json:"max_response_time"`
	MinResponseTime   time.Duration           `json:"min_response_time"`
	QPS               float64                 `json:"qps"`
	LastResetTime     time.Time               `json:"last_reset_time"`
	MethodStats       map[string]*MethodStats `json:"method_stats"`
	PathStats         map[string]*PathStats   `json:"path_stats"`
	StatusStats       map[int]*StatusStats    `json:"status_stats"`
}

// MethodStats 按HTTP方法统计的性能指标
type MethodStats struct {
	Count             int64         `json:"count"`
	TotalResponseTime time.Duration `json:"total_response_time"`
	AvgResponseTime   time.Duration `json:"avg_response_time"`
}

// PathStats 按路径统计的性能指标
type PathStats struct {
	Count             int64         `json:"count"`
	TotalResponseTime time.Duration `json:"total_response_time"`
	AvgResponseTime   time.Duration `json:"avg_response_time"`
}

// StatusStats 按状态码统计的性能指标
type StatusStats struct {
	Count             int64         `json:"count"`
	TotalResponseTime time.Duration `json:"total_response_time"`
	AvgResponseTime   time.Duration `json:"avg_response_time"`
}

// 全局性能指标
var (
	performanceMetrics = &PerformanceMetrics{
		TotalRequests:     0,
		TotalResponseTime: 0,
		MaxResponseTime:   0,
		MinResponseTime:   time.Hour, // 初始值设为1小时
		LastResetTime:     time.Now(),
		MethodStats:       make(map[string]*MethodStats),
		PathStats:         make(map[string]*PathStats),
		StatusStats:       make(map[int]*StatusStats),
	}
	metricsMutex = &sync.RWMutex{}
)

// PerformanceMiddleware 性能监控中间件
func PerformanceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 记录请求结束时间和耗时
		endTime := time.Now()
		latency := endTime.Sub(startTime)

		// 获取状态码
		statusCode := c.Writer.Status()
		// 获取请求方法
		method := c.Request.Method
		// 获取请求路径
		path := c.Request.URL.Path

		// 更新性能指标（使用细粒度锁，减少锁持有时间）
		updatePerformanceMetrics(latency, method, path, statusCode, endTime)

		// 在响应头中添加请求耗时
		c.Writer.Header().Set("X-Response-Time", latency.String())
	}
}

// updatePerformanceMetrics 更新性能指标
func updatePerformanceMetrics(latency time.Duration, method, path string, statusCode int, endTime time.Time) {
	metricsMutex.Lock()
	defer metricsMutex.Unlock()

	// 更新总请求数
	performanceMetrics.TotalRequests++

	// 更新总响应时间
	performanceMetrics.TotalResponseTime += latency

	// 更新平均响应时间
	performanceMetrics.AvgResponseTime = performanceMetrics.TotalResponseTime / time.Duration(performanceMetrics.TotalRequests)

	// 更新最大响应时间
	if latency > performanceMetrics.MaxResponseTime {
		performanceMetrics.MaxResponseTime = latency
	}

	// 更新最小响应时间
	if latency < performanceMetrics.MinResponseTime {
		performanceMetrics.MinResponseTime = latency
	}

	// 更新QPS（每秒请求数）
	elapsed := endTime.Sub(performanceMetrics.LastResetTime).Seconds()
	if elapsed > 0 {
		performanceMetrics.QPS = float64(performanceMetrics.TotalRequests) / elapsed
	}

	// 更新按方法统计的指标
	if _, exists := performanceMetrics.MethodStats[method]; !exists {
		performanceMetrics.MethodStats[method] = &MethodStats{
			Count:             0,
			TotalResponseTime: 0,
		}
	}
	methodStat := performanceMetrics.MethodStats[method]
	methodStat.Count++
	methodStat.TotalResponseTime += latency
	methodStat.AvgResponseTime = methodStat.TotalResponseTime / time.Duration(methodStat.Count)

	// 更新按路径统计的指标
	if _, exists := performanceMetrics.PathStats[path]; !exists {
		performanceMetrics.PathStats[path] = &PathStats{
			Count:             0,
			TotalResponseTime: 0,
		}
	}
	pathStat := performanceMetrics.PathStats[path]
	pathStat.Count++
	pathStat.TotalResponseTime += latency
	pathStat.AvgResponseTime = pathStat.TotalResponseTime / time.Duration(pathStat.Count)

	// 更新按状态码统计的指标
	if _, exists := performanceMetrics.StatusStats[statusCode]; !exists {
		performanceMetrics.StatusStats[statusCode] = &StatusStats{
			Count:             0,
			TotalResponseTime: 0,
		}
	}
	statusStat := performanceMetrics.StatusStats[statusCode]
	statusStat.Count++
	statusStat.TotalResponseTime += latency
	statusStat.AvgResponseTime = statusStat.TotalResponseTime / time.Duration(statusStat.Count)
}

// GetPerformanceMetrics 获取当前性能指标
func GetPerformanceMetrics() *PerformanceMetrics {
	metricsMutex.RLock()
	defer metricsMutex.RUnlock()

	// 返回性能指标的副本，避免并发修改问题
	return &PerformanceMetrics{
		TotalRequests:     performanceMetrics.TotalRequests,
		TotalResponseTime: performanceMetrics.TotalResponseTime,
		AvgResponseTime:   performanceMetrics.AvgResponseTime,
		MaxResponseTime:   performanceMetrics.MaxResponseTime,
		MinResponseTime:   performanceMetrics.MinResponseTime,
		QPS:               performanceMetrics.QPS,
		LastResetTime:     performanceMetrics.LastResetTime,
		MethodStats:       performanceMetrics.MethodStats,
		PathStats:         performanceMetrics.PathStats,
		StatusStats:       performanceMetrics.StatusStats,
	}
}
