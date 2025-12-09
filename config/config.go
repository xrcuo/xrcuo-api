package config

import (
	_ "embed"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
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
		Level      string `yaml:"level"`
		File       string `yaml:"file"`
		MaxSize    int    `yaml:"max_size"`    // 单个日志文件最大大小（MB）
		MaxBackups int    `yaml:"max_backups"` // 保留的日志文件数量
		MaxAge     int    `yaml:"max_age"`     // 日志文件保留天数
	} `yaml:"log"`

	RandomImage struct {
		LocalEnabled bool   `yaml:"local_enabled"` // 是否启用本地图片
		LocalPath    string `yaml:"local_path"`    // 本地图片目录路径
	} `yaml:"random_image"`
}

var Conf *Config

// 生成配置文件
func genConfig() error {
	configPath := getConfigPath()
	logrus.Debugf("正在生成配置文件: %s", configPath)
	return os.WriteFile(configPath, []byte(defConfig), 0644)
}

// 获取配置文件路径，支持从环境变量CONFIG_FILE_PATH指定
func getConfigPath() string {
	if path := os.Getenv("CONFIG_FILE_PATH"); path != "" {
		return path
	}
	return "config.yaml"
}

// 解析配置文件
func Parse() {
	// 使用环境变量或默认路径
	configPath := getConfigPath()

	logrus.Debugf("正在解析配置文件: %s", configPath)

	// 检查文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		err = genConfig()
		if err != nil {
			logrus.Fatalf("无法生成设置文件: %s, 请确认是否给足系统权限", configPath)
		}
		logrus.Warnf("未检测到 %s，已自动生成，请配置并重新启动", configPath)
		logrus.Warn("将于 5 秒后退出...")
		os.Exit(1)
	}

	// 读取配置文件
	content, err := os.ReadFile(configPath)
	if err != nil {
		logrus.Fatalf("读取配置文件失败: %v", err)
	}

	// 解析配置文件
	Conf = &Config{}
	err = yaml.Unmarshal(content, Conf)
	if err != nil {
		logrus.Fatal("解析 config.yaml 失败，请检查格式、内容是否输入正确")
	}

	// 验证配置
	validateConfig()

	// 应用日志级别配置
	setLogLevel()
}

// validateConfig 验证配置的有效性
func validateConfig() {
	// 验证Gin运行模式
	validModes := map[string]bool{
		"debug":   true,
		"release": true,
		"test":    true,
	}
	if !validModes[Conf.Server.Mode] {
		logrus.Warnf("无效的Gin模式: %s, 使用默认模式: debug", Conf.Server.Mode)
		Conf.Server.Mode = "debug"
	}

	// 向后兼容：处理旧版本配置
	// 检查是否存在旧的 db_path 配置
	if Conf.IP2Region.V4DBPath == "" && Conf.IP2Region.V6DBPath == "" {
		// 尝试读取旧的 config 文件中的 db_path 字段
		// 这里需要使用反射或手动解析，因为 yaml.Unmarshal 不会将不存在的字段设置为默认值
		// 我们将使用默认值来处理
		logrus.Warn("检测到旧版本配置格式，将使用默认配置")
		// 使用默认值
		Conf.IP2Region.V4DBPath = "./ip2region_v4.xdb"
		Conf.IP2Region.V6DBPath = "./ip2region_v6.xdb"
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
	if !validLogLevels[Conf.Log.Level] {
		logrus.Warnf("无效的日志级别: %s, 使用默认级别: info", Conf.Log.Level)
		Conf.Log.Level = "info"
	}

	// 验证日志文件大小
	if Conf.Log.MaxSize <= 0 {
		logrus.Warnf("无效的日志文件大小: %d, 使用默认值: 10 MB", Conf.Log.MaxSize)
		Conf.Log.MaxSize = 10
	}

	// 验证日志文件保留数量
	if Conf.Log.MaxBackups <= 0 {
		logrus.Warnf("无效的日志文件保留数量: %d, 使用默认值: 5", Conf.Log.MaxBackups)
		Conf.Log.MaxBackups = 5
	}

	// 验证日志文件保留天数
	if Conf.Log.MaxAge <= 0 {
		logrus.Warnf("无效的日志文件保留天数: %d, 使用默认值: 7", Conf.Log.MaxAge)
		Conf.Log.MaxAge = 7
	}

	logrus.Debug("配置验证完成")
}

// GetServerPort 获取服务器端口
func GetServerPort() string {
	if Conf == nil || Conf.Server.Port == "" {
		return ":8080"
	}
	return Conf.Server.Port
}

// GetServerMode 获取Gin运行模式
func GetServerMode() string {
	if Conf == nil || Conf.Server.Mode == "" {
		return "debug"
	}
	return Conf.Server.Mode
}

// IsJSONFormatEnabled 获取是否启用JSON格式化
func IsJSONFormatEnabled() bool {
	if Conf == nil {
		return false
	}
	return Conf.Server.JSONFormat.Enabled
}

// GetDatabasePath 获取数据库文件路径
func GetDatabasePath() string {
	if Conf == nil || Conf.Database.Path == "" {
		return "./stats.db"
	}
	return Conf.Database.Path
}

// GetMaxOpenConns 获取最大打开连接数
func GetMaxOpenConns() int {
	if Conf == nil || Conf.Database.MaxOpenConns <= 0 {
		return 10
	}
	return Conf.Database.MaxOpenConns
}

// GetMaxIdleConns 获取最大空闲连接数
func GetMaxIdleConns() int {
	if Conf == nil || Conf.Database.MaxIdleConns <= 0 {
		return 5
	}
	return Conf.Database.MaxIdleConns
}

// GetIP2RegionV4DBPath 获取IP2Region IPv4数据库路径
func GetIP2RegionV4DBPath() string {
	if Conf == nil || Conf.IP2Region.V4DBPath == "" {
		return "./ip2region_v4.xdb"
	}
	return Conf.IP2Region.V4DBPath
}

// GetIP2RegionV6DBPath 获取IP2Region IPv6数据库路径
func GetIP2RegionV6DBPath() string {
	if Conf == nil || Conf.IP2Region.V6DBPath == "" {
		return "./ip2region_v6.xdb"
	}
	return Conf.IP2Region.V6DBPath
}

// GetLogLevel 获取日志级别
func GetLogLevel() string {
	if Conf == nil || Conf.Log.Level == "" {
		return "info"
	}
	return Conf.Log.Level
}

// 根据配置设置日志级别和文件输出
func setLogLevel() {
	// 设置日志级别
	levelStr := GetLogLevel()
	level, err := logrus.ParseLevel(levelStr)
	if err != nil {
		logrus.Warnf("无效的日志级别: %s, 使用默认级别: info", levelStr)
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)

	// 设置日志格式
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// 如果配置了日志文件路径，则添加文件输出
	if Conf.Log.File != "" {
		// 确保日志目录存在
		logDir := filepath.Dir(Conf.Log.File)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			logrus.Warnf("创建日志目录失败: %v, 仅输出到控制台", err)
			return
		}

		// 配置日志文件输出，使用 lumberjack 处理日志滚动
		logrus.SetOutput(&lumberjack.Logger{
			Filename:   Conf.Log.File,
			MaxSize:    Conf.Log.MaxSize,
			MaxBackups: Conf.Log.MaxBackups,
			MaxAge:     Conf.Log.MaxAge,
		})

		logrus.Infof("日志文件已配置: %s", Conf.Log.File)
	}

	logrus.Debugf("日志级别已设置为: %s", level)
}
