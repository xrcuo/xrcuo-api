package random

import (
	"github.com/gin-gonic/gin"
)

// RandomPlugin Random插件实现
var RandomPlugin = &randomPlugin{}

// randomPlugin Random插件结构体
type randomPlugin struct{}

// Name 返回插件名称
func (p *randomPlugin) Name() string {
	return "random"
}

// RegisterRouter 注册随机图片插件路由
func (p *randomPlugin) RegisterRouter(group *gin.RouterGroup) {
	{
		// 获取随机图片（重定向到图片URL）
		group.GET("/random/image", GetRandomImageHandler)
		// 获取随机图片信息（返回JSON格式）
		group.GET("/random/image/info", GetRandomImageInfoHandler)
	}
}

// RegisterRouter 注册随机图片插件的路由（兼容旧的注册方式）
func RegisterRouter(router *gin.RouterGroup) {
	RandomPlugin.RegisterRouter(router)
}
