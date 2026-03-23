package mcpe

import (
	"time"
)

// Response 响应结构
type Response struct {
	Code int         `json:"code"` // 状态码
	Msg  string      `json:"msg"`  // 消息
	Data interface{} `json:"data"` // 数据
	Took string      `json:"took"` // 耗时
}

// Data 响应数据结构
type Data struct {
	ServerIP   string    `json:"server_ip"`   // 服务器IP
	Port       int       `json:"port"`        // 服务器端口
	Online     int       `json:"online"`      // 在线人数
	MaxPlayers int       `json:"max_players"` // 最大玩家数
	Version    string    `json:"version"`     // 服务器版本
	Motd       string    `json:"motd"`        // 服务器描述
	PingTime   string    `json:"ping_time"`   // 延迟
	Time       time.Time `json:"time"`        // 查询时间
}
