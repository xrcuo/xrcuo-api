package config

import "github.com/lionsoul2014/ip2region/binding/golang/xdb"

// 全局配置项
var (
	IP2RegionDBPath = "./ip2region_v4.xdb" // ip2region数据库路径
	IPVersion       = xdb.IPv4             // IP版本（IPv4）
	ServerPort      = ":8080"              // 服务监听端口
)
