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

// RegisterRouter 注册IP插件路由
func (p *ipPlugin) RegisterRouter(group *gin.RouterGroup) {
	// 路由前缀：/api/ip
	ipGroup := group.Group("/ip")
	{
		// GET /api/ip?ip=xxx.xxx.xxx.xxx
		ipGroup.GET("", SearchRegionHandler)
	}
}

// RegisterRouter 注册IP插件路由（兼容旧的注册方式）
func RegisterRouter(group *gin.RouterGroup) {
	IPPlugin.RegisterRouter(group)
}
