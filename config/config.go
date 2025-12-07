package config

import (
	_ "embed"
	"os"
	"strings"

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
	sb := strings.Builder{}
	sb.WriteString(defConfig)
	err := os.WriteFile("config.yaml", []byte(sb.String()), 0644)
	if err != nil {
		return err
	}
	return nil
}

// 解析配置文件
func Parse() {
	// 使用绝对路径确保配置文件位置正确
	configPath := "config.yaml"

	// 检查文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		err = genConfig()
		if err != nil {
			panic("无法生成设置文件: config.yaml, 请确认是否给足系统权限")
		}
		logrus.Warn("未检测到 config.yaml，已自动于同目录生成，请配置并重新启动")
		logrus.Warn("将于 5 秒后退出...")
		os.Exit(-1)
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
