package client

import (
	"github.com/gin-gonic/gin"
)

// ClientPlugin Client插件实现
var ClientPlugin = &clientPlugin{}

// clientPlugin Client插件结构体
type clientPlugin struct{}

// Name 返回插件名称
func (p *clientPlugin) Name() string {
	return "client"
}

// Init 初始化插件
func (p *clientPlugin) Init() error {
	// Client插件初始化逻辑
	return nil
}

// RegisterRouter 注册客户端信息插件路由
func (p *clientPlugin) RegisterRouter(group *gin.RouterGroup) {
	// 注册客户端信息API路由，路径为/api/client
	group.GET("/client", GetClientInfoHandler)
}

// Cleanup 清理插件资源
func (p *clientPlugin) Cleanup() error {
	// Client插件清理逻辑
	return nil
}

// RegisterRouter 注册客户端信息API路由（兼容旧的注册方式）
func RegisterRouter(routerGroup *gin.RouterGroup) {
	ClientPlugin.RegisterRouter(routerGroup)
}
