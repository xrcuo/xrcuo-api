package ip

import "github.com/gin-gonic/gin"

// IPPlugin IP插件实现
var IPPlugin = &ipPlugin{}

// ipPlugin IP插件结构体
type ipPlugin struct{}

// Name 返回插件名称
func (p *ipPlugin) Name() string {
	return "ip"
}

// Init 初始化插件
func (p *ipPlugin) Init() error {
	// IP插件初始化逻辑
	return nil
}

// RegisterRouter 注册IP插件路由
func (p *ipPlugin) RegisterRouter(group *gin.RouterGroup) {
	// 路由前缀：/api/ip
	ipGroup := group.Group("/ip")
	{
		// GET /api/ip?ip=xxx.xxx.xxx.xxx
		ipGroup.GET("", SearchRegionHandler)
	}
}

// Cleanup 清理插件资源
func (p *ipPlugin) Cleanup() error {
	// IP插件清理逻辑
	return nil
}

// RegisterRouter 注册IP插件路由（兼容旧的注册方式）
func RegisterRouter(group *gin.RouterGroup) {
	IPPlugin.RegisterRouter(group)
}
