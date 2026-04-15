package config

import (
	_ "embed"
	"log"
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

type Config struct {
	Server struct {
		Port       string `yaml:"port"`
		Mode       string `yaml:"mode"`
		JSONFormat struct {
			Enabled bool `yaml:"enabled"`
		} `yaml:"json_format"`
	} `yaml:"server"`

	Database struct {
		Type         string `yaml:"type"`
		Path         string `yaml:"path"`
		Host         string `yaml:"host"`
		Port         int    `yaml:"port"`
		User         string `yaml:"user"`
		Password     string `yaml:"password"`
		DBName       string `yaml:"dbname"`
		MaxOpenConns int    `yaml:"max_open_conns"`
		MaxIdleConns int    `yaml:"max_idle_conns"`
	} `yaml:"database"`

	IP2Region struct {
		V4DBPath string `yaml:"v4_db_path"`
		V6DBPath string `yaml:"v6_db_path"`
	} `yaml:"ip2region"`

	Log struct {
		Level            string `yaml:"level"`
		File             string `yaml:"file"`
		ConsoleOutput    bool   `yaml:"console_output"`
		RequestLog       bool   `yaml:"request_log"`
		MaxSize          int    `yaml:"max_size"`
		MaxBackups       int    `yaml:"max_backups"`
		MaxAge           int    `yaml:"max_age"`
		Compress         bool   `yaml:"compress"`
		NewFileOnStartup bool   `yaml:"new_file_on_startup"`
	} `yaml:"log"`
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
	debounceTimer   *time.Timer
	debounceMutex   sync.Mutex
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
			stopChan: make(chan struct{}, 1),
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

func (cm *ConfigManager) genConfig() error {
	logrus.Debugf("正在生成配置文件: %s", cm.configPath)
	return os.WriteFile(cm.configPath, []byte(defConfig), 0644)
}

func (cm *ConfigManager) getConfigPath() string {
	if path := os.Getenv("CONFIG_FILE_PATH"); path != "" {
		return path
	}
	return "config.yaml"
}

func Parse() {
	GetInstance().ParseConfig()
}

func (cm *ConfigManager) RegisterUpdateCallback(callback ConfigUpdateCallback) {
	cm.callbacksMutex.Lock()
	defer cm.callbacksMutex.Unlock()
	cm.updateCallbacks = append(cm.updateCallbacks, callback)
}

func (cm *ConfigManager) executeUpdateCallbacks(config *Config) {
	cm.callbacksMutex.Lock()
	callbacks := make([]ConfigUpdateCallback, len(cm.updateCallbacks))
	copy(callbacks, cm.updateCallbacks)
	cm.callbacksMutex.Unlock()

	for _, callback := range callbacks {
		callback(config)
	}
}

func (cm *ConfigManager) ParseConfig() {
	cm.configPath = cm.getConfigPath()

	logrus.Debugf("正在解析配置文件: %s", cm.configPath)

	if _, err := os.Stat(cm.configPath); os.IsNotExist(err) {
		err = cm.genConfig()
		if err != nil {
			logrus.Fatalf("无法生成设置文件: %s, 请确认是否给足系统权限", cm.configPath)
		}
		logrus.Warnf("未检测到 %s，已自动生成，请配置并重新启动", cm.configPath)
		logrus.Warn("将于 5 秒后退出...")
		os.Exit(1)
	}

	content, err := os.ReadFile(cm.configPath)
	if err != nil {
		logrus.Fatalf("读取配置文件失败: %v", err)
	}

	config := &Config{}
	err = yaml.Unmarshal(content, config)
	if err != nil {
		logrus.Fatal("解析 config.yaml 失败，请检查格式、内容是否输入正确")
	}

	cm.validateConfig(config)

	isUpdate := cm.config != nil
	cm.SetConfig(config)

	if isUpdate {
		logrus.Info("正在应用更新后的配置...")
		gin.SetMode(cm.GetConfig().Server.Mode)
		logrus.Infof("Gin模式已更新为: %s", cm.GetConfig().Server.Mode)
		cm.executeUpdateCallbacks(config)
		log.Println("配置更新应用完成")
	}
}

func (cm *ConfigManager) validateConfig(config *Config) {
	validModes := map[string]bool{"debug": true, "release": true, "test": true}
	if !validModes[config.Server.Mode] {
		logrus.Warnf("无效的Gin模式: %s, 使用默认模式: debug", config.Server.Mode)
		config.Server.Mode = "debug"
	}

	if config.IP2Region.V4DBPath == "" && config.IP2Region.V6DBPath == "" {
		logrus.Warn("检测到旧版本配置格式，将使用默认配置")
		config.IP2Region.V4DBPath = "./ip2region_v4.xdb"
		config.IP2Region.V6DBPath = "./ip2region_v6.xdb"
	}

	validLogLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true, "fatal": true, "panic": true}
	if !validLogLevels[config.Log.Level] {
		logrus.Warnf("无效的日志级别: %s, 使用默认级别: info", config.Log.Level)
		config.Log.Level = "info"
	}

	if config.Log.MaxSize <= 0 {
		logrus.Warnf("无效的日志文件大小: %d, 使用默认值: 10 MB", config.Log.MaxSize)
		config.Log.MaxSize = 10
	}

	if config.Log.MaxBackups <= 0 {
		logrus.Warnf("无效的日志文件保留数量: %d, 使用默认值: 5", config.Log.MaxBackups)
		config.Log.MaxBackups = 5
	}

	if config.Log.MaxAge <= 0 {
		logrus.Warnf("无效的日志文件保留天数: %d, 使用默认值: 7", config.Log.MaxAge)
		config.Log.MaxAge = 7
	}

	logrus.Debug("配置验证完成")
}

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

	if err := cm.watcher.Add(cm.configPath); err != nil {
		logrus.Errorf("添加配置文件到监听器失败: %v", err)
		cm.watcher.Close()
		return
	}

	cm.isWatching = true

	go func() {
		defer func() {
			cm.watcher.Close()
			cm.isWatching = false
			cm.debounceMutex.Lock()
			if cm.debounceTimer != nil {
				cm.debounceTimer.Stop()
			}
			cm.debounceMutex.Unlock()
		}()

		for {
			select {
			case event, ok := <-cm.watcher.Events:
				if !ok {
					return
				}

				if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
					logrus.Info("配置文件发生变化，重新加载配置")
					cm.debounceMutex.Lock()
					if cm.debounceTimer != nil {
						cm.debounceTimer.Stop()
					}
					cm.debounceTimer = time.AfterFunc(100*time.Millisecond, func() {
						cm.ParseConfig()
					})
					cm.debounceMutex.Unlock()
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

func (cm *ConfigManager) StopWatching() {
	if !cm.isWatching {
		return
	}

	select {
	case cm.stopChan <- struct{}{}:
	default:
	}
	cm.isWatching = false
	logrus.Info("配置文件监听已停止")
}

func GetServerPort() string {
	cm := GetInstance()
	config := cm.GetConfig()
	if config == nil || config.Server.Port == "" {
		return ":8080"
	}
	return config.Server.Port
}

func GetServerMode() string {
	cm := GetInstance()
	config := cm.GetConfig()
	if config == nil || config.Server.Mode == "" {
		return "debug"
	}
	return config.Server.Mode
}

func IsJSONFormatEnabled() bool {
	cm := GetInstance()
	config := cm.GetConfig()
	if config == nil {
		return false
	}
	return config.Server.JSONFormat.Enabled
}

func GetDatabasePath() string {
	cm := GetInstance()
	config := cm.GetConfig()
	if config == nil || config.Database.Path == "" {
		return "./stats.db"
	}
	return config.Database.Path
}

func GetMaxOpenConns() int {
	cm := GetInstance()
	config := cm.GetConfig()
	if config == nil || config.Database.MaxOpenConns <= 0 {
		return 10
	}
	return config.Database.MaxOpenConns
}

func GetMaxIdleConns() int {
	cm := GetInstance()
	config := cm.GetConfig()
	if config == nil || config.Database.MaxIdleConns <= 0 {
		return 5
	}
	return config.Database.MaxIdleConns
}

func GetIP2RegionV4DBPath() string {
	cm := GetInstance()
	config := cm.GetConfig()
	if config == nil || config.IP2Region.V4DBPath == "" {
		return "./ip2region_v4.xdb"
	}
	return config.IP2Region.V4DBPath
}

func GetIP2RegionV6DBPath() string {
	cm := GetInstance()
	config := cm.GetConfig()
	if config == nil || config.IP2Region.V6DBPath == "" {
		return "./ip2region_v6.xdb"
	}
	return config.IP2Region.V6DBPath
}

func GetLogLevel() string {
	cm := GetInstance()
	config := cm.GetConfig()
	if config == nil || config.Log.Level == "" {
		return "info"
	}
	return config.Log.Level
}

func GetDatabaseType() string {
	cm := GetInstance()
	config := cm.GetConfig()
	if config == nil || config.Database.Type == "" {
		return "sqlite"
	}
	return config.Database.Type
}

func GetDatabaseHost() string {
	cm := GetInstance()
	config := cm.GetConfig()
	if config == nil || config.Database.Host == "" {
		return "localhost"
	}
	return config.Database.Host
}

func GetDatabasePort() int {
	cm := GetInstance()
	config := cm.GetConfig()
	if config == nil || config.Database.Port <= 0 {
		dbType := GetDatabaseType()
		switch dbType {
		case "mysql":
			return 3306
		case "postgresql":
			return 5432
		}
	}
	return config.Database.Port
}

func GetDatabaseUser() string {
	cm := GetInstance()
	config := cm.GetConfig()
	if config == nil {
		return ""
	}
	return config.Database.User
}

func GetDatabasePassword() string {
	cm := GetInstance()
	config := cm.GetConfig()
	if config == nil {
		return ""
	}
	return config.Database.Password
}

func GetDatabaseName() string {
	cm := GetInstance()
	config := cm.GetConfig()
	if config == nil || config.Database.DBName == "" {
		return "xrcuo_api"
	}
	return config.Database.DBName
}
