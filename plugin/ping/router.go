package ping

import "github.com/gin-gonic/gin"

// PingPlugin Ping插件实现
var PingPlugin = &pingPlugin{}

// pingPlugin Ping插件结构体
type pingPlugin struct{}

// Name 返回插件名称
func (p *pingPlugin) Name() string {
	return "ping"
}

// RegisterRouter 注册Ping插件路由
func (p *pingPlugin) RegisterRouter(group *gin.RouterGroup) {
	// 路由前缀：/api/ping
	pingGroup := group.Group("/ping")
	{
		// GET /api/ping?target=www.baidu.com&count=3
		pingGroup.GET("", PingHandler)
	}
}

// RegisterRouter 注册Ping插件路由（兼容旧的注册方式）
func RegisterRouter(group *gin.RouterGroup) {
	PingPlugin.RegisterRouter(group)
}
