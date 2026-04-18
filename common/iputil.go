package common

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

// IsPrivateIP 判断是否为内网IP、回环地址或保留地址
func IsPrivateIP(ip string) bool {
	ipAddr := net.ParseIP(ip)
	if ipAddr == nil {
		return false
	}
	// 检查是否为回环地址（127.0.0.1 或 ::1）
	if ipAddr.IsLoopback() {
		return true
	}
	// 检查是否为私有内网IP
	if ipAddr.IsPrivate() {
		return true
	}
	// 额外检查其他保留地址范围
	// 198.18.0.0/15 (用于基准测试)
	if ip4 := ipAddr.To4(); ip4 != nil {
		if (ip4[0] == 198 && ip4[1] >= 18 && ip4[1] <= 19) ||
			// 100.64.0.0/10 (运营商级NAT)
			(ip4[0] == 100 && (ip4[1]&0xc0) == 0x40) ||
			// 192.0.0.0/24 (IETF协议分配)
			(ip4[0] == 192 && ip4[1] == 0 && ip4[2] == 0) ||
			// 192.0.2.0/24 (文档和示例)
			(ip4[0] == 192 && ip4[1] == 0 && ip4[2] == 2) ||
			// 198.51.100.0/24 (文档和示例)
			(ip4[0] == 198 && ip4[1] == 51 && ip4[2] == 100) ||
			// 203.0.113.0/24 (文档和示例)
			(ip4[0] == 203 && ip4[1] == 0 && ip4[2] == 113) ||
			// 224.0.0.0/4 (多播地址)
			(ip4[0] >= 224 && ip4[0] <= 239) ||
			// 255.255.255.255 (广播地址)
			(ip4[0] == 255 && ip4[1] == 255 && ip4[2] == 255 && ip4[3] == 255) {
			return true
		}
	}
	return false
}

// ResolveTarget 解析目标（域名→IP，IP直接返回）
func ResolveTarget(target string) (string, error) {
	// 先判断是否为IP地址
	if net.ParseIP(target) != nil {
		return target, nil
	}

	// 域名解析（优先IPv4）
	ips, err := net.DefaultResolver.LookupIP(context.Background(), "ip4", target)
	if err != nil {
		return "", fmt.Errorf("域名解析失败：%v", err)
	}
	if len(ips) == 0 {
		return "", fmt.Errorf("未解析到IPv4地址")
	}

	return ips[0].String(), nil
}

// GetLocalIPs 获取本地所有网络接口的IP地址
func GetLocalIPs() ([]string, error) {
	var ips []string

	// 获取所有网络接口
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range interfaces {
		// 跳过非活动接口和回环接口
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		// 获取接口的地址
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			// 解析IP地址
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			// 只添加IPv4地址
			if ip != nil && ip.To4() != nil {
				ips = append(ips, ip.String())
			}
		}
	}

	return ips, nil
}

// GetPublicIP 获取当前外网IP（使用国内API）
func GetPublicIP() (string, error) {
	// 国内轻量IP查询服务
	services := []string{
		"http://ip.3322.net",
		"http://myip.ipip.net",
		"http://ifconfig.co/ip",
		"http://checkip.amazonaws.com",
		"http://ident.me",
	}

	client := &http.Client{
		Timeout: 8 * time.Second,
	}

	for _, service := range services {
		req, err := http.NewRequest("GET", service, nil)
		if err != nil {
			continue
		}

		// 设置用户代理，避免被某些服务拒绝
		req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; xrcuo-api)")

		resp, err := client.Do(req)
		if err != nil {
			continue
		}

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			continue
		}

		// 清理响应内容
		ip := strings.TrimSpace(string(body))

		// 特殊处理 ipip.net 的响应格式
		if service == "http://myip.ipip.net" {
			// ipip.net 返回格式: 当前 IP: 123.45.67.89 来自: 北京市 联通
			parts := strings.Split(ip, " ")
			if len(parts) >= 3 {
				ip = parts[2]
			}
		}

		// 验证IP格式
		if net.ParseIP(ip) != nil {
			return ip, nil
		}
	}

	return "", fmt.Errorf("无法获取外网IP")
}
