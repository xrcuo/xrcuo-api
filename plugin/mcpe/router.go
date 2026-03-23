package mcpe

import (
	"github.com/gin-gonic/gin"
)

// MCPEPlugin MCPE插件实例
var MCPEPlugin = &mcpePlugin{}

// mcpePlugin 插件结构体
type mcpePlugin struct{}

// Name 返回插件名称
func (p *mcpePlugin) Name() string {
	return "mcpe"
}

// Init 初始化插件
func (p *mcpePlugin) Init() error {
	return nil
}

// RegisterRouter 注册插件路由
func (p *mcpePlugin) RegisterRouter(group *gin.RouterGroup) {
	mcpeGroup := group.Group("/mcpe")
	{
		mcpeGroup.GET("/status", MCPEHandler)
	}
}

// Cleanup 清理插件资源
func (p *mcpePlugin) Cleanup() error {
	return nil
}
