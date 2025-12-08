package ipify

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xrcuo/xrcuo-api/plugin/client"
)

// GetIPHandler 获取客户端IP地址处理函数
func GetIPHandler(c *gin.Context) {
	// 获取客户端真实IP
	clientIP := client.GetRealIP(c)

	// 直接返回IP地址，不使用JSON格式
	c.String(http.StatusOK, clientIP)
}
