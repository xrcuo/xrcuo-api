package common

import (
	"fmt"
	"strings"
	"sync"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
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

// 预加载的IP2Region查询器
var preloadedSearcher *xdb.Searcher

// 公共对象池（复用ip2region查询器，减少资源占用）
var regionSearcherPool = sync.Pool{
	New: func() interface{} {
		return preloadedSearcher
	},
}

// InitIP2Region 预加载IP2Region数据库
func InitIP2Region() error {
	dbPath := config.GetIP2RegionDBPath()
	ipVersion := config.GetIPVersion()

	logrus.Infof("开始预加载IP2Region数据库: %s, IP版本: %v", dbPath, ipVersion)

	// 预加载数据库到内存
	searcher, err := xdb.NewWithFileOnly(ipVersion, dbPath)
	if err != nil {
		return fmt.Errorf("IP2Region数据库预加载失败: %v", err)
	}

	preloadedSearcher = searcher
	logrus.Info("IP2Region数据库预加载成功")
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

	// 从对象池获取查询器
	searcher := regionSearcherPool.Get()
	if searcher == nil {
		return RegionParts{}, fmt.Errorf("地区查询器初始化失败")
	}
	defer regionSearcherPool.Put(searcher) // 归还对象池

	// 执行查询
	regionRaw, err := searcher.(*xdb.Searcher).SearchByStr(ip)
	if err != nil {
		return RegionParts{}, fmt.Errorf("IP查询失败：%v", err)
	}

	// 解析地区字符串（适配4段格式：国家|省份|城市|ISP）
	return parseRegionRaw(regionRaw), nil
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
