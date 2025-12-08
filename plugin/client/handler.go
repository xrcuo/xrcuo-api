package client

import (
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xrcuo/xrcuo-api/common"
)

// GetClientInfoHandler 获取客户端信息处理函数
func GetClientInfoHandler(c *gin.Context) {
	startTime := time.Now()
	response := &Response{Code: 200, Msg: "请求成功"}

	// 统一响应出口（确保took字段必赋值）
	defer func() {
		response.Took = time.Since(startTime).String()
		c.JSON(http.StatusOK, response)
	}()

	// 1. 获取客户端真实IP
	clientIP := GetRealIP(c)

	// 2. 调用公共工具查询IP地区信息
	regionParts, err := common.GetRegionByIP(clientIP)

	// 处理IP查询失败的情况，仍然返回客户端信息，只是地区信息为空
	if err != nil {
		regionParts = common.RegionParts{
			Country:  "",
			Province: "",
			City:     "",
			Isp:      "",
		}
	}

	// 3. 解析User-Agent获取操作系统和浏览器信息
	userAgent := c.Request.UserAgent()
	os, browser, browserVersion := parseUserAgent(userAgent)

	// 4. 构造响应数据
	locationParts := []string{regionParts.Country, regionParts.Province, regionParts.City}
	location := common.JoinNonEmpty(locationParts, "")
	area := common.JoinNonEmpty(append(locationParts, regionParts.Isp), "")

	response.Data = &Data{
		IP:             clientIP,
		Location:       location,
		ISP:            regionParts.Isp,
		Area:           area,
		OS:             os,
		Browser:        browser,
		BrowserVersion: browserVersion,
	}
}

// parseUserAgent 解析User-Agent字符串，获取操作系统和浏览器信息
func parseUserAgent(ua string) (os string, browser string, version string) {
	// 默认为未知
	os = "Unknown"
	browser = "Unknown"
	version = "Unknown"

	// 简化的User-Agent解析，实际项目中可以使用成熟的库如github.com/mssola/user_agent

	// 解析操作系统
	osPatterns := map[string]string{
		"Windows NT 10.0": "Windows 10",
		"Windows NT 6.3":  "Windows 8.1",
		"Windows NT 6.2":  "Windows 8",
		"Windows NT 6.1":  "Windows 7",
		"Windows NT 6.0":  "Windows Vista",
		"Windows NT 5.1":  "Windows XP",
		"Mac OS X":        "macOS",
		"iPhone OS":       "iOS",
		"Android":         "Android",
		"Linux":           "Linux",
	}

	for pattern, name := range osPatterns {
		if strings.Contains(ua, pattern) {
			os = name
			break
		}
	}

	// 解析浏览器 - 按优先级顺序匹配，确保Chrome优先于Safari
	browserPatterns := []struct {
		Name    string
		Pattern *regexp.Regexp
	}{
		{"Google Chrome", regexp.MustCompile(`Chrome/(\d+\.\d+\.\d+\.\d+)`)},
		{"Microsoft Edge", regexp.MustCompile(`Edg/(\d+\.\d+\.\d+\.\d+)`)},
		{"Mozilla Firefox", regexp.MustCompile(`Firefox/(\d+\.\d+)`)},
		{"Opera", regexp.MustCompile(`OPR/(\d+\.\d+\.\d+\.\d+)`)},
		{"Safari", regexp.MustCompile(`Version/(\d+\.\d+)\s+Safari`)},
	}

	for _, bp := range browserPatterns {
		matches := bp.Pattern.FindStringSubmatch(ua)
		if len(matches) > 1 {
			browser = bp.Name
			version = matches[1]
			break
		}
	}

	return
}

// GetRealIP 获取客户端真实IP，处理代理情况
func GetRealIP(c *gin.Context) string {
	// 优先从X-Real-IP头获取
	if ip := c.GetHeader("X-Real-IP"); ip != "" {
		return ip
	}

	// 然后从X-Forwarded-For头获取，取第一个IP
	if ip := c.GetHeader("X-Forwarded-For"); ip != "" {
		ips := strings.Split(ip, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// 最后使用Gin的ClientIP方法
	return c.ClientIP()
}

// getPublicIP 调用外部服务获取真实公网IP（备用方案）
func getPublicIP() (string, error) {
	// 使用多个备选服务，提高可靠性
	services := []string{
		"https://api.ipify.org",
		"https://ipinfo.io/ip",
		"https://icanhazip.com",
	}

	for _, service := range services {
		resp, err := http.Get(service)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			ipBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				continue
			}
			ip := strings.TrimSpace(string(ipBytes))
			if ip != "" {
				return ip, nil
			}
		}
	}

	// 如果所有服务都失败，返回默认的示例IP
	return "112.51.209.121", nil
}
