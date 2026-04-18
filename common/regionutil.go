package common

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
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

	// 检查并下载 IPv4 数据库文件
	if err := checkAndDownloadDB(v4DBPath, "IPv4"); err != nil {
		return fmt.Errorf("检查/下载IPv4数据库失败: %v", err)
	}

	// 检查并下载 IPv6 数据库文件（IPv6是可选的，失败不影响启动）
	if err := checkAndDownloadDB(v6DBPath, "IPv6"); err != nil {
		logrus.Warnf("检查/下载IPv6数据库失败，将仅使用IPv4: %v", err)
	}

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

	logrus.Debugf("IP查询原始结果: IP=%s, Raw=%s", ip, regionRaw)

	// 解析地区字符串（适配4段格式：国家|省份|城市|ISP）
	result := parseRegionRaw(regionRaw)
	logrus.Debugf("IP查询解析结果: %+v", result)

	return result, nil
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

	// 定义一个辅助函数，保留更多数据并过滤掉无效值
	keepField := func(field string) string {
		field = strings.TrimSpace(field)
		if field == "0" || field == "" || strings.EqualFold(field, "reserved") {
			return ""
		}
		return field
	}

	switch len(parts) {
	case 5:
		result.Isp = keepField(parts[3])
		result.City = keepField(parts[2])
		result.Province = keepField(parts[1])
		result.Country = keepField(parts[0])
	case 4:
		result.Isp = keepField(parts[3])
		result.City = keepField(parts[2])
		result.Province = keepField(parts[1])
		result.Country = keepField(parts[0])
	case 3:
		result.City = keepField(parts[2])
		result.Province = keepField(parts[1])
		result.Country = keepField(parts[0])
	case 2:
		result.Province = keepField(parts[1])
		result.Country = keepField(parts[0])
	case 1:
		result.Country = keepField(parts[0])
	}

	return result
}

// JoinNonEmpty 合并非空字符串（忽略空值）
func JoinNonEmpty(strs []string, sep string) string {
	var result []string
	for _, s := range strs {
		if s != "" {
			result = append(result, s)
		}
	}
	return strings.Join(result, sep)
}

// downloadFile 下载文件到指定路径
func downloadFile(url, filePath string) error {
	logrus.Infof("正在下载文件: %s -> %s", url, filePath)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("下载请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("下载失败，HTTP状态码: %d", resp.StatusCode)
	}

	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	file, fileErr := os.Create(filePath)
	if fileErr != nil {
		return fmt.Errorf("创建文件失败: %v", fileErr)
	}
	defer file.Close()

	_, copyErr := io.Copy(file, resp.Body)
	if copyErr != nil {
		return fmt.Errorf("写入文件失败: %v", copyErr)
	}

	logrus.Infof("文件下载成功: %s", filePath)
	return nil
}

// checkAndDownloadDB 检查数据库文件是否存在，不存在则自动下载
func checkAndDownloadDB(dbPath, dbType string) error {
	if _, err := os.Stat(dbPath); err == nil {
		logrus.Infof("%s数据库文件已存在: %s", dbType, dbPath)
		return nil
	}

	logrus.Warnf("%s数据库文件不存在，准备自动下载: %s", dbType, dbPath)

	var downloadURL string
	switch dbType {
	case "IPv4":
		downloadURL = "https://github.com/lionsoul2014/ip2region/raw/master/data/ip2region_v4.xdb"
	case "IPv6":
		downloadURL = "https://github.com/lionsoul2014/ip2region/raw/master/data/ip2region_v6.xdb"
	default:
		return fmt.Errorf("未知的数据库类型: %s", dbType)
	}

	if err := downloadFile(downloadURL, dbPath); err != nil {
		return fmt.Errorf("下载%s数据库失败: %v", dbType, err)
	}

	return nil
}
