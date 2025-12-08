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

// RegisterRouter 注册客户端信息插件路由
func (p *clientPlugin) RegisterRouter(group *gin.RouterGroup) {
	// 注册客户端信息API路由，路径为/api/client
	group.GET("/client", GetClientInfoHandler)
}

// RegisterRouter 注册客户端信息API路由（兼容旧的注册方式）
func RegisterRouter(routerGroup *gin.RouterGroup) {
	ClientPlugin.RegisterRouter(routerGroup)
}
