package models

import (
	"time"
)

// Stats 存储API调用统计信息
type Stats struct {
	TotalCalls      int64            `json:"total_calls"`       // 总调用次数
	DailyCalls      int64            `json:"daily_calls"`       // 今日调用次数
	HourlyCalls     int64            `json:"hourly_calls"`      // 每小时调用次数
	MethodCalls     map[string]int64 `json:"method_calls"`      // 按HTTP方法统计
	PathCalls       map[string]int64 `json:"path_calls"`        // 按API路径统计
	IPCalls         map[string]int64 `json:"ip_calls"`          // 按IP统计
	LastResetTime   time.Time        `json:"last_reset_time"`   // 上次重置时间
	LastCallDetails []*CallDetail    `json:"last_call_details"` // 最近调用详情
}

// CallDetail 存储单个API调用的详细信息
type CallDetail struct {
	Path       string    `json:"path"`        // 请求路径
	Method     string    `json:"method"`      // 请求方法
	IP         string    `json:"ip"`          // 请求IP
	Timestamp  time.Time `json:"timestamp"`   // 请求时间
	StatusCode int       `json:"status_code"` // 响应状态码
}
