package common

import (
	"github.com/gin-gonic/gin"
)

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
