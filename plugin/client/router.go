package client

import (
	"github.com/gin-gonic/gin"
)

// RegisterRouter 注册客户端信息API路由
func RegisterRouter(routerGroup *gin.RouterGroup) {
	// 注册客户端信息API路由，路径为/api/client
	routerGroup.GET("/client", GetClientInfoHandler)
}
