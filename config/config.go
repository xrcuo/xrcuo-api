package config

import (
	_ "embed"
	"os"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

//go:embed default_config.yaml
var defConfig string

// Config 应用程序配置结构体
type Config struct {
	Server struct {
		Port       string `yaml:"port"`
		Mode       string `yaml:"mode"` // Gin运行模式（debug, release, test）
		JSONFormat struct {
			Enabled bool `yaml:"enabled"` // 是否启用格式化JSON响应
		} `yaml:"json_format"`
	} `yaml:"server"`

	Database struct {
		Path         string `yaml:"path"`           // SQLite数据库文件路径
		MaxOpenConns int    `yaml:"max_open_conns"` // 最大打开连接数
		MaxIdleConns int    `yaml:"max_idle_conns"` // 最大空闲连接数
	} `yaml:"database"`

	IP2Region struct {
		V4DBPath string `yaml:"v4_db_path"` // IPv4数据库文件路径
		V6DBPath string `yaml:"v6_db_path"` // IPv6数据库文件路径
	} `yaml:"ip2region"`

	Log struct {
		Level         string `yaml:"level"`
		File          string `yaml:"file"`
		ConsoleOutput bool   `yaml:"console_output"` // 是否输出到控制台
		RequestLog    bool   `yaml:"request_log"`    // 是否输出请求日志
		MaxSize       int    `yaml:"max_size"`       // 单个日志文件最大大小（MB）
		MaxBackups    int    `yaml:"max_backups"`    // 保留的日志文件数量
		MaxAge        int    `yaml:"max_age"`        // 日志文件保留天数
	} `yaml:"log"`

	RandomImage struct {
		LocalEnabled bool   `yaml:"local_enabled"` // 是否启用本地图片
		LocalPath    string `yaml:"local_path"`    // 本地图片目录路径
	} `yaml:"random_image"`
}

// ConfigUpdateCallback 配置更新回调函数类型
type ConfigUpdateCallback func(*Config)

// ConfigManager 配置管理器单例
type ConfigManager struct {
	config          *Config
	configPath      string
	mutex           sync.RWMutex
	watcher         *fsnotify.Watcher
	stopChan        chan struct{}
	isWatching      bool
	updateCallbacks []ConfigUpdateCallback
	callbacksMutex  sync.Mutex
}

// 全局配置管理器实例
var (
	instance *ConfigManager
	once     sync.Once
)

// GetInstance 获取配置管理器单例
func GetInstance() *ConfigManager {
	once.Do(func() {
		instance = &ConfigManager{
			stopChan: make(chan struct{}),
		}
	})
	return instance
}

// GetConfig 获取当前配置
func (cm *ConfigManager) GetConfig() *Config {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	return cm.config
}

// SetConfig 设置配置
func (cm *ConfigManager) SetConfig(config *Config) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.config = config
}

// 生成配置文件
func (cm *ConfigManager) genConfig() error {
	logrus.Debugf("正在生成配置文件: %s", cm.configPath)
	return os.WriteFile(cm.configPath, []byte(defConfig), 0644)
}

// 获取配置文件路径，支持从环境变量CONFIG_FILE_PATH指定
func (cm *ConfigManager) getConfigPath() string {
	if path := os.Getenv("CONFIG_FILE_PATH"); path != "" {
		return path
	}
	return "config.yaml"
}

// Parse 解析配置文件
func Parse() {
	cm := GetInstance()
	cm.ParseConfig()
}

// RegisterUpdateCallback 注册配置更新回调
func (cm *ConfigManager) RegisterUpdateCallback(callback ConfigUpdateCallback) {
	cm.callbacksMutex.Lock()
	defer cm.callbacksMutex.Unlock()
	cm.updateCallbacks = append(cm.updateCallbacks, callback)
}

// executeUpdateCallbacks 执行所有配置更新回调
func (cm *ConfigManager) executeUpdateCallbacks(config *Config) {
	cm.callbacksMutex.Lock()
	callbacks := make([]ConfigUpdateCallback, len(cm.updateCallbacks))
	copy(callbacks, cm.updateCallbacks)
	cm.callbacksMutex.Unlock()

	for _, callback := range callbacks {
		callback(config)
	}
}

// ParseConfig 解析配置文件
func (cm *ConfigManager) ParseConfig() {
	// 使用环境变量或默认路径
	cm.configPath = cm.getConfigPath()

	logrus.Debugf("正在解析配置文件: %s", cm.configPath)

	// 检查文件是否存在
	if _, err := os.Stat(cm.configPath); os.IsNotExist(err) {
		err = cm.genConfig()
		if err != nil {
			logrus.Fatalf("无法生成设置文件: %s, 请确认是否给足系统权限", cm.configPath)
		}
		logrus.Warnf("未检测到 %s，已自动生成，请配置并重新启动", cm.configPath)
		logrus.Warn("将于 5 秒后退出...")
		os.Exit(1)
	}

	// 读取配置文件
	content, err := os.ReadFile(cm.configPath)
	if err != nil {
		logrus.Fatalf("读取配置文件失败: %v", err)
	}

	// 解析配置文件
	config := &Config{}
	err = yaml.Unmarshal(content, config)
	if err != nil {
		logrus.Fatal("解析 config.yaml 失败，请检查格式、内容是否输入正确")
	}

	// 验证配置
	cm.validateConfig(config)

	// 检查是否是配置更新
	isUpdate := cm.config != nil

	// 设置配置
	cm.SetConfig(config)

	// 如果是配置更新，则执行更新逻辑
	if isUpdate {
		logrus.Info("正在应用更新后的配置...")

		// 更新Gin模式
		gin.SetMode(cm.GetConfig().Server.Mode)
		logrus.Infof("Gin模式已更新为: %s", cm.GetConfig().Server.Mode)

		// 执行所有配置更新回调
		cm.executeUpdateCallbacks(config)

		logrus.Info("配置更新应用完成")
	}
}

// validateConfig 验证配置的有效性
func (cm *ConfigManager) validateConfig(config *Config) {
	// 验证Gin运行模式
	validModes := map[string]bool{
		"debug":   true,
		"release": true,
		"test":    true,
	}
	if !validModes[config.Server.Mode] {
		logrus.Warnf("无效的Gin模式: %s, 使用默认模式: debug", config.Server.Mode)
		config.Server.Mode = "debug"
	}

	// 向后兼容：处理旧版本配置
	// 检查是否存在旧的 db_path 配置
	if config.IP2Region.V4DBPath == "" && config.IP2Region.V6DBPath == "" {
		// 尝试读取旧的 config 文件中的 db_path 字段
		// 这里需要使用反射或手动解析，因为 yaml.Unmarshal 不会将不存在的字段设置为默认值
		// 我们将使用默认值来处理
		logrus.Warn("检测到旧版本配置格式，将使用默认配置")
		// 使用默认值
		config.IP2Region.V4DBPath = "./ip2region_v4.xdb"
		config.IP2Region.V6DBPath = "./ip2region_v6.xdb"
	}

	// 验证日志级别
	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
		"fatal": true,
		"panic": true,
	}
	if !validLogLevels[config.Log.Level] {
		logrus.Warnf("无效的日志级别: %s, 使用默认级别: info", config.Log.Level)
		config.Log.Level = "info"
	}

	// 验证日志文件大小
	if config.Log.MaxSize <= 0 {
		logrus.Warnf("无效的日志文件大小: %d, 使用默认值: 10 MB", config.Log.MaxSize)
		config.Log.MaxSize = 10
	}

	// 验证日志文件保留数量
	if config.Log.MaxBackups <= 0 {
		logrus.Warnf("无效的日志文件保留数量: %d, 使用默认值: 5", config.Log.MaxBackups)
		config.Log.MaxBackups = 5
	}

	// 验证日志文件保留天数
	if config.Log.MaxAge <= 0 {
		logrus.Warnf("无效的日志文件保留天数: %d, 使用默认值: 7", config.Log.MaxAge)
		config.Log.MaxAge = 7
	}

	logrus.Debug("配置验证完成")
}

// WatchConfig 监听配置文件变化
func (cm *ConfigManager) WatchConfig() {
	if cm.isWatching {
		return
	}

	var err error
	cm.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		logrus.Errorf("创建配置文件监听器失败: %v", err)
		return
	}

	// 添加配置文件到监听器
	if err := cm.watcher.Add(cm.configPath); err != nil {
		logrus.Errorf("添加配置文件到监听器失败: %v", err)
		cm.watcher.Close()
		return
	}

	cm.isWatching = true

	// 启动监听协程
	go func() {
		defer func() {
			cm.watcher.Close()
			cm.isWatching = false
		}()

		for {
			select {
			case event, ok := <-cm.watcher.Events:
				if !ok {
					return
				}

				// 只处理写入和创建事件
				if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
					logrus.Info("配置文件发生变化，重新加载配置")
					// 延迟处理，避免文件被锁定
					time.Sleep(100 * time.Millisecond)
					cm.ParseConfig()
				}
			case err, ok := <-cm.watcher.Errors:
				if !ok {
					return
				}
				logrus.Errorf("配置文件监听错误: %v", err)
			case <-cm.stopChan:
				return
			}
		}
	}()

	logrus.Info("配置文件监听已启动")
}

// StopWatching 停止监听配置文件
func (cm *ConfigManager) StopWatching() {
	if !cm.isWatching {
		return
	}

	cm.stopChan <- struct{}{}
	cm.isWatching = false
	logrus.Info("配置文件监听已停止")
}

// GetServerPort 获取服务器端口
func GetServerPort() string {
	cm := GetInstance()
	config := cm.GetConfig()
	if config == nil || config.Server.Port == "" {
		return ":8080"
	}
	return config.Server.Port
}

// GetServerMode 获取Gin运行模式
func GetServerMode() string {
	cm := GetInstance()
	config := cm.GetConfig()
	if config == nil || config.Server.Mode == "" {
		return "debug"
	}
	return config.Server.Mode
}

// IsJSONFormatEnabled 获取是否启用JSON格式化
func IsJSONFormatEnabled() bool {
	cm := GetInstance()
	config := cm.GetConfig()
	if config == nil {
		return false
	}
	return config.Server.JSONFormat.Enabled
}

// GetDatabasePath 获取数据库文件路径
func GetDatabasePath() string {
	cm := GetInstance()
	config := cm.GetConfig()
	if config == nil || config.Database.Path == "" {
		return "./stats.db"
	}
	return config.Database.Path
}

// GetMaxOpenConns 获取最大打开连接数
func GetMaxOpenConns() int {
	cm := GetInstance()
	config := cm.GetConfig()
	if config == nil || config.Database.MaxOpenConns <= 0 {
		return 10
	}
	return config.Database.MaxOpenConns
}

// GetMaxIdleConns 获取最大空闲连接数
func GetMaxIdleConns() int {
	cm := GetInstance()
	config := cm.GetConfig()
	if config == nil || config.Database.MaxIdleConns <= 0 {
		return 5
	}
	return config.Database.MaxIdleConns
}

// GetIP2RegionV4DBPath 获取IP2Region IPv4数据库路径
func GetIP2RegionV4DBPath() string {
	cm := GetInstance()
	config := cm.GetConfig()
	if config == nil || config.IP2Region.V4DBPath == "" {
		return "./ip2region_v4.xdb"
	}
	return config.IP2Region.V4DBPath
}

// GetIP2RegionV6DBPath 获取IP2Region IPv6数据库路径
func GetIP2RegionV6DBPath() string {
	cm := GetInstance()
	config := cm.GetConfig()
	if config == nil || config.IP2Region.V6DBPath == "" {
		return "./ip2region_v6.xdb"
	}
	return config.IP2Region.V6DBPath
}

// GetLogLevel 获取日志级别
func GetLogLevel() string {
	cm := GetInstance()
	config := cm.GetConfig()
	if config == nil || config.Log.Level == "" {
		return "info"
	}
	return config.Log.Level
}
