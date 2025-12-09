package common

import (
	"fmt"
	"strings"

	"github.com/lionsoul2014/ip2region/binding/golang/service"
	"github.com/sirupsen/logrus"
	"github.com/xrcuo/xrcuo-api/config"
)

// RegionParts 地区结构化数据
type RegionParts struct {
	Country  string // 国家
	Province string // 省份
	City     string // 城市
	Isp      string // 运营商
}

// 全局ip2region服务
var ip2regionService *service.Ip2Region

// InitIP2Region 初始化IP2Region服务
func InitIP2Region() error {
	v4DBPath := config.GetIP2RegionV4DBPath()
	v6DBPath := config.GetIP2RegionV6DBPath()

	logrus.Infof("开始初始化IP2Region服务: IPv4路径: %s, IPv6路径: %s", v4DBPath, v6DBPath)

	// 创建v4配置：指定缓存策略和v4的xdb文件路径
	v4Config, err := service.NewV4Config(service.VIndexCache, v4DBPath, 20)
	if err != nil {
		return fmt.Errorf("创建IPv4配置失败: %v", err)
	}

	// 尝试创建v6配置，如果失败则只使用v4配置
	v6Config, err := service.NewV6Config(service.VIndexCache, v6DBPath, 20)
	if err != nil {
		logrus.Warnf("创建IPv6配置失败，将只使用IPv4配置: %v", err)
		// 通过配置创建Ip2Region查询服务（只使用v4配置）
		ip2regionService, err = service.NewIp2Region(v4Config, nil)
	} else {
		// 通过配置创建Ip2Region查询服务（同时使用v4和v6配置）
		ip2regionService, err = service.NewIp2Region(v4Config, v6Config)
	}

	if err != nil {
		return fmt.Errorf("创建IP2Region服务失败: %v", err)
	}

	logrus.Info("IP2Region服务初始化成功")
	return nil
}

// GetRegionByIP 根据IP获取地区信息（支持内网IP识别）
func GetRegionByIP(ip string) (RegionParts, error) {
	// 内网IP直接返回固定结果
	if IsPrivateIP(ip) {
		return RegionParts{
			Country:  "内网IP",
			Province: "",
			City:     "",
			Isp:      "",
		}, nil
	}

	// 执行查询
	regionRaw, err := ip2regionService.SearchByStr(ip)
	if err != nil {
		return RegionParts{}, fmt.Errorf("IP查询失败：%v", err)
	}

	// 解析地区字符串（适配4段格式：国家|省份|城市|ISP）
	return parseRegionRaw(regionRaw), nil
}

// CloseIP2Region 关闭IP2Region服务
func CloseIP2Region() {
	if ip2regionService != nil {
		ip2regionService.Close()
		logrus.Info("IP2Region服务已关闭")
	}
}

// parseRegionRaw 解析ip2region原始返回值
func parseRegionRaw(regionRaw string) RegionParts {
	parts := strings.Split(regionRaw, "|")
	result := RegionParts{}

	switch len(parts) {
	case 4:
		result.Isp = parseEmptyField(parts[3])
		result.City = parseEmptyField(parts[2])
		result.Province = parseEmptyField(parts[1])
		result.Country = parseEmptyField(parts[0])
	case 3:
		result.City = parseEmptyField(parts[2])
		result.Province = parseEmptyField(parts[1])
		result.Country = parseEmptyField(parts[0])
	case 2:
		result.Province = parseEmptyField(parts[1])
		result.Country = parseEmptyField(parts[0])
	case 1:
		result.Country = parseEmptyField(parts[0])
	}

	return result
}

// parseEmptyField 处理空字段（"0"或空字符串转为""）
func parseEmptyField(field string) string {
	if field == "0" || field == "" || field == "未知" {
		return ""
	}
	return field
}

// JoinNonEmpty 合并非空字符串（忽略空值，无分隔符）
func JoinNonEmpty(strs []string, sep string) string {
	var result []string
	for _, s := range strs {
		if s != "" {
			result = append(result, s)
		}
	}
	return strings.Join(result, sep)
}
