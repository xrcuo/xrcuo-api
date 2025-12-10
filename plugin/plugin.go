package plugin

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xrcuo/xrcuo-api/plugin/api_key"
	"github.com/xrcuo/xrcuo-api/plugin/client"
	"github.com/xrcuo/xrcuo-api/plugin/ip"
	"github.com/xrcuo/xrcuo-api/plugin/ipify"
	"github.com/xrcuo/xrcuo-api/plugin/ping"
	"github.com/xrcuo/xrcuo-api/plugin/random"
)

// Plugin 插件接口
type Plugin interface {
	// Name 返回插件名称
	Name() string
	// Init 初始化插件
	Init() error
	// RegisterRouter 注册插件路由
	RegisterRouter(group *gin.RouterGroup)
	// Cleanup 清理插件资源
	Cleanup() error
}

// PluginInfo 插件信息
type PluginInfo struct {
	Name    string
	Enabled bool
}

// PluginManager 插件管理器
type PluginManager struct {
	plugins     []Plugin
	initialized bool
	pluginInfos map[string]*PluginInfo
}

// NewPluginManager 创建新的插件管理器
func NewPluginManager() *PluginManager {
	return &PluginManager{
		plugins:     make([]Plugin, 0),
		pluginInfos: make(map[string]*PluginInfo),
	}
}

// Register 注册插件
func (pm *PluginManager) Register(plugin Plugin) {
	name := plugin.Name()
	if _, exists := pm.pluginInfos[name]; exists {
		logrus.Warnf("插件 %s 已注册，跳过重复注册", name)
		return
	}

	pm.plugins = append(pm.plugins, plugin)
	pm.pluginInfos[name] = &PluginInfo{
		Name:    name,
		Enabled: true,
	}

	logrus.Infof("插件 %s 已注册", name)
}

// InitAll 初始化所有插件
func (pm *PluginManager) InitAll() error {
	if pm.initialized {
		return nil
	}

	for _, plugin := range pm.plugins {
		if err := plugin.Init(); err != nil {
			logrus.Errorf("初始化插件 %s 失败：%v", plugin.Name(), err)
			return err
		}
		logrus.Infof("插件 %s 初始化成功", plugin.Name())
	}

	pm.initialized = true
	return nil
}

// RegisterAll 注册所有插件到指定路由组
func (pm *PluginManager) RegisterAll(group *gin.RouterGroup) {
	for _, plugin := range pm.plugins {
		plugin.RegisterRouter(group)
		logrus.Infof("插件 %s 路由注册成功", plugin.Name())
	}
}

// CleanupAll 清理所有插件资源
func (pm *PluginManager) CleanupAll() {
	for _, plugin := range pm.plugins {
		if err := plugin.Cleanup(); err != nil {
			logrus.Errorf("清理插件 %s 资源失败：%v", plugin.Name(), err)
			continue
		}
		logrus.Infof("插件 %s 资源清理成功", plugin.Name())
	}

	pm.initialized = false
}

// GetPlugins 获取所有注册的插件
func (pm *PluginManager) GetPlugins() []Plugin {
	return pm.plugins
}

// GetPluginInfo 获取插件信息
func (pm *PluginManager) GetPluginInfo(name string) (*PluginInfo, bool) {
	info, exists := pm.pluginInfos[name]
	return info, exists
}

// EnablePlugin 启用插件
func (pm *PluginManager) EnablePlugin(name string) bool {
	info, exists := pm.pluginInfos[name]
	if !exists {
		return false
	}

	info.Enabled = true
	logrus.Infof("插件 %s 已启用", name)
	return true
}

// DisablePlugin 禁用插件
func (pm *PluginManager) DisablePlugin(name string) bool {
	info, exists := pm.pluginInfos[name]
	if !exists {
		return false
	}

	info.Enabled = false
	logrus.Infof("插件 %s 已禁用", name)
	return true
}

// RegisterBuiltinPlugins 注册所有内置插件
func (pm *PluginManager) RegisterBuiltinPlugins() {
	pm.Register(ip.IPPlugin)
	pm.Register(ping.PingPlugin)
	pm.Register(random.RandomPlugin)
	pm.Register(client.ClientPlugin)
	pm.Register(ipify.IpifyPlugin)
}

// RegisterAPIRouter 注册API密钥管理路由
func RegisterAPIRouter(r *gin.RouterGroup) {
	api_key.RegisterRouter(r)
}
