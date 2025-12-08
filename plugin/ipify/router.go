package ipify

import (
	"github.com/gin-gonic/gin"
)

// IpifyPlugin Ipify插件实现
var IpifyPlugin = &ipifyPlugin{}

// ipifyPlugin Ipify插件结构体
type ipifyPlugin struct{}

// Name 返回插件名称
func (p *ipifyPlugin) Name() string {
	return "ipify"
}

// RegisterRouter 注册IPify插件路由
func (p *ipifyPlugin) RegisterRouter(group *gin.RouterGroup) {
	// 注册IP获取API路由，路径为/api/ipify
	group.GET("/ipify", GetIPHandler)
}

// RegisterRouter 注册IP获取API路由（兼容旧的注册方式）
func RegisterRouter(routerGroup *gin.RouterGroup) {
	IpifyPlugin.RegisterRouter(routerGroup)
}
