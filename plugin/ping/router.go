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

// Init 初始化插件
func (p *pingPlugin) Init() error {
	// Ping插件初始化逻辑
	return nil
}

// RegisterRouter 注册Ping插件路由
func (p *pingPlugin) RegisterRouter(group *gin.RouterGroup) {
	// 路由前缀：/api/ping
	pingGroup := group.Group("/ping")
	{
		// GET /api/ping?target=xxx&count=3
		pingGroup.GET("", PingHandler)
	}
}

// Cleanup 清理插件资源
func (p *pingPlugin) Cleanup() error {
	// Ping插件清理逻辑
	return nil
}

// RegisterRouter 注册Ping插件路由（兼容旧的注册方式）
func RegisterRouter(group *gin.RouterGroup) {
	PingPlugin.RegisterRouter(group)
}
