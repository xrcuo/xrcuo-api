package common

import (
	"context"
	"fmt"
	"net"
)

// IsPrivateIP 判断是否为内网IP
func IsPrivateIP(ip string) bool {
	ipAddr := net.ParseIP(ip)
	if ipAddr == nil {
		return false
	}
	return ipAddr.IsPrivate() // 内置方法判断内网IP（10.0.0.0/8、172.16.0.0/12、192.168.0.0/16）
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
