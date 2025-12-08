package common

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

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
		
		// 记录请求日志
		logrus.WithFields(logrus.Fields{
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"status":     c.Writer.Status(),
			"client_ip":  c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
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
		// 优化跨域配置，添加更多安全头
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		c.Writer.Header().Set("Access-Control-Max-Age", "3600") // 预检请求结果缓存1小时
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		
		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
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
