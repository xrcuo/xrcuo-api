package plugin

import "github.com/gin-gonic/gin"

// Plugin 插件接口
type Plugin interface {
	// Name 返回插件名称
	Name() string
	// RegisterRouter 注册插件路由
	RegisterRouter(group *gin.RouterGroup)
}

// PluginManager 插件管理器
type PluginManager struct {
	plugins []Plugin
}

// NewPluginManager 创建新的插件管理器
func NewPluginManager() *PluginManager {
	return &PluginManager{
		plugins: make([]Plugin, 0),
	}
}

// Register 注册插件
func (pm *PluginManager) Register(plugin Plugin) {
	pm.plugins = append(pm.plugins, plugin)
}

// RegisterAll 注册所有插件到指定路由组
func (pm *PluginManager) RegisterAll(group *gin.RouterGroup) {
	for _, plugin := range pm.plugins {
		plugin.RegisterRouter(group)
	}
}

// GetPlugins 获取所有注册的插件
func (pm *PluginManager) GetPlugins() []Plugin {
	return pm.plugins
}
