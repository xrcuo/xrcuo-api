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

// Init 初始化插件
func (p *ipifyPlugin) Init() error {
	// Ipify插件初始化逻辑
	return nil
}

// RegisterRouter 注册IPify插件路由
func (p *ipifyPlugin) RegisterRouter(group *gin.RouterGroup) {
	// 注册IP获取API路由，路径为/api/ipify
	group.GET("/ipify", GetIPHandler)
}

// Cleanup 清理插件资源
func (p *ipifyPlugin) Cleanup() error {
	// Ipify插件清理逻辑
	return nil
}

// RegisterRouter 注册IP获取API路由（兼容旧的注册方式）
func RegisterRouter(routerGroup *gin.RouterGroup) {
	IpifyPlugin.RegisterRouter(routerGroup)
}
