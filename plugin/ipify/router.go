package ipify

import (
	"github.com/gin-gonic/gin"
)

// RegisterRouter 注册IP获取API路由
func RegisterRouter(routerGroup *gin.RouterGroup) {
	// 注册IP获取API路由，路径为/api/ipify
	routerGroup.GET("/ipify", GetIPHandler)
}
