package ip

import "github.com/gin-gonic/gin"

// RegisterRouter 注册IP插件路由
func RegisterRouter(group *gin.RouterGroup) {
	// 路由前缀：/api/ip
	ipGroup := group.Group("/ip")
	{
		// GET /api/ip?ip=xxx.xxx.xxx.xxx
		ipGroup.GET("", SearchRegionHandler)
	}
}
