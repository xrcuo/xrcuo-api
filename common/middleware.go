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
			"size":       c.Writer.Size(),
		}).Info("API请求")
	}
}

// CORSMiddleware 跨域请求中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

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
		// 处理请求前的操作
		path := c.Request.URL.Path
		method := c.Request.Method
		clientIP := c.ClientIP()

		// 处理请求
		c.Next()

		// 处理请求后的操作
		statusCode := c.Writer.Status()

		// 记录调用信息
		if GlobalStats != nil {
			GlobalStats.RecordCall(path, method, clientIP, statusCode)
		}
	}
}
