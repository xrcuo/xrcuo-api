package config

import (
	_ "embed"
	"os"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

//go:embed default_config.yaml
var defConfig string

// Config 应用程序配置结构体
type Config struct {
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`

	IP2Region struct {
		DBPath    string `yaml:"db_path"`
		IPVersion string `yaml:"ip_version"` // ipv4 or ipv6
	} `yaml:"ip2region"`

	Log struct {
		Level string `yaml:"level"`
	} `yaml:"log"`
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

	// 应用日志级别配置
	setLogLevel()
}

// GetServerPort 获取服务器端口
func GetServerPort() string {
	if Conf == nil || Conf.Server.Port == "" {
		return ":8080"
	}
	return Conf.Server.Port
}

// GetIP2RegionDBPath 获取IP2Region数据库路径
func GetIP2RegionDBPath() string {
	if Conf == nil || Conf.IP2Region.DBPath == "" {
		return "./ip2region_v4.xdb"
	}
	return Conf.IP2Region.DBPath
}

// GetIPVersion 获取IP版本
func GetIPVersion() *xdb.Version {
	if Conf == nil {
		return xdb.IPv4
	}
	// 根据配置的IP版本返回对应的常量
	if Conf.IP2Region.IPVersion == "ipv6" {
		return xdb.IPv6
	}
	return xdb.IPv4 // 默认IPv4
}

// GetLogLevel 获取日志级别
func GetLogLevel() string {
	if Conf == nil || Conf.Log.Level == "" {
		return "info"
	}
	return Conf.Log.Level
}

// 根据配置设置日志级别
func setLogLevel() {
	levelStr := GetLogLevel()
	level, err := logrus.ParseLevel(levelStr)
	if err != nil {
		logrus.Warnf("无效的日志级别: %s, 使用默认级别: info", levelStr)
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.Debugf("日志级别已设置为: %s", level)
}
