package ping

import "github.com/gin-gonic/gin"

// RegisterRouter 注册Ping插件路由
func RegisterRouter(group *gin.RouterGroup) {
	// 路由前缀：/api/ping
	pingGroup := group.Group("/ping")
	{
		// GET /api/ping?target=xxx&count=3&timeout=5
		pingGroup.GET("", PingHandler)
	}
}
