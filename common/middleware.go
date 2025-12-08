package common

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xrcuo/xrcuo-api/db"
	"github.com/xrcuo/xrcuo-api/models"
)

// apiKeyCacheItem API密钥缓存项
type apiKeyCacheItem struct {
	Key        *models.APIKey
	ExpireTime time.Time
}

// apiKeyCache API密钥缓存
type apiKeyCache struct {
	items          map[string]*apiKeyCacheItem
	mutex          sync.RWMutex
	expireDuration time.Duration
}

// 全局API密钥缓存实例
var apiKeyCacheInstance = &apiKeyCache{
	items:          make(map[string]*apiKeyCacheItem),
	expireDuration: 5 * time.Minute, // 缓存5分钟
}

// Get 获取缓存的API密钥
func (c *apiKeyCache) Get(key string) *models.APIKey {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return nil
	}

	// 检查是否过期
	if time.Now().After(item.ExpireTime) {
		// 过期后异步删除
		go c.Delete(key)
		return nil
	}

	return item.Key
}

// Set 设置API密钥缓存
func (c *apiKeyCache) Set(key string, apiKey *models.APIKey) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.items[key] = &apiKeyCacheItem{
		Key:        apiKey,
		ExpireTime: time.Now().Add(c.expireDuration),
	}
}

// Delete 删除API密钥缓存
func (c *apiKeyCache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.items, key)
}

// Clear 清空API密钥缓存
func (c *apiKeyCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.items = make(map[string]*apiKeyCacheItem)
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

// CORSMiddleware 跨域请求中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 安全的CORS配置，避免使用通配符
		origin := c.GetHeader("Origin")
		if origin != "" {
			// 允许特定来源（生产环境中应替换为实际域名）
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			// 如果没有Origin头，允许所有来源
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		c.Writer.Header().Set("Access-Control-Max-Age", "3600") // 预检请求结果缓存1小时
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")

		// 安全头：防止点击劫持
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		// 安全头：防止XSS攻击
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		// 安全头：防止MIME类型嗅探
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")

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
		keyInfo := apiKeyCacheInstance.Get(apiKey)
		var err error
		if keyInfo == nil {
			// 从数据库获取
			keyInfo, err = db.GetAPIKeyByKey(apiKey)
			if err != nil {
				ErrorResponse(c, http.StatusUnauthorized, 401, "无效的API密钥")
				c.Abort()
				return
			}
			// 存入缓存
			apiKeyCacheInstance.Set(apiKey, keyInfo)
		}

		// 检查API密钥是否已达到使用上限
		if !keyInfo.IsPermanent && keyInfo.CurrentUsage >= keyInfo.MaxUsage {
			ErrorResponse(c, http.StatusForbidden, 403, "API密钥已达到使用上限")
			c.Abort()
			return
		}

		// 更新API密钥使用次数
		if err := db.UpdateAPIKeyUsage(apiKey); err != nil {
			ErrorResponse(c, http.StatusInternalServerError, 500, "更新API密钥使用次数失败")
			c.Abort()
			return
		}

		// 更新缓存中的使用次数
		keyInfo.CurrentUsage++
		apiKeyCacheInstance.Set(apiKey, keyInfo)

		// 将API密钥信息存储到上下文
		c.Set("api_key", keyInfo)

		// 继续处理请求
		c.Next()
	}
}

// rateLimitItem 速率限制项
type rateLimitItem struct {
	Count      int
	LastAccess time.Time
}

// rateLimiter 速率限制器
type rateLimiter struct {
	items           map[string]*rateLimitItem
	mutex           sync.RWMutex
	limit           int           // 时间窗口内的最大请求数
	window          time.Duration // 时间窗口
	cleanupInterval time.Duration // 清理过期项的时间间隔
}

// 全局速率限制器实例
var globalRateLimiter = &rateLimiter{
	items:           make(map[string]*rateLimitItem),
	limit:           100,         // 100个请求
	window:          time.Minute, // 1分钟窗口
	cleanupInterval: time.Hour,   // 每小时清理一次过期项
}

// init 初始化速率限制器
func init() {
	// 启动定期清理过期项的任务
	go func() {
		ticker := time.NewTicker(globalRateLimiter.cleanupInterval)
		defer ticker.Stop()

		for range ticker.C {
			globalRateLimiter.cleanupExpiredItems()
		}
	}()
}

// Allow 检查是否允许请求
func (rl *rateLimiter) Allow(key string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	item, exists := rl.items[key]
	now := time.Now()

	// 如果不存在，或者已经过期，初始化新的限制项
	if !exists || now.Sub(item.LastAccess) > rl.window {
		rl.items[key] = &rateLimitItem{
			Count:      1,
			LastAccess: now,
		}
		return true
	}

	// 如果请求数超过限制，拒绝请求
	if item.Count >= rl.limit {
		return false
	}

	// 否则，增加请求数并更新最后访问时间
	item.Count++
	item.LastAccess = now
	return true
}

// cleanupExpiredItems 清理过期的速率限制项
func (rl *rateLimiter) cleanupExpiredItems() {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	for key, item := range rl.items {
		if now.Sub(item.LastAccess) > rl.window {
			delete(rl.items, key)
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
