package common

import (
	"sync"
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
	mu              sync.RWMutex     // 读写锁，保证并发安全
}

// CallDetail 存储单个API调用的详细信息
type CallDetail struct {
	Path       string    `json:"path"`        // 请求路径
	Method     string    `json:"method"`      // 请求方法
	IP         string    `json:"ip"`          // 请求IP
	Timestamp  time.Time `json:"timestamp"`   // 请求时间
	StatusCode int       `json:"status_code"` // 响应状态码
}

// 全局统计实例
var GlobalStats *Stats

// InitStats 初始化统计信息
func InitStats() {
	GlobalStats = &Stats{
		MethodCalls:     make(map[string]int64),
		PathCalls:       make(map[string]int64),
		IPCalls:         make(map[string]int64),
		LastResetTime:   time.Now(),
		LastCallDetails: make([]*CallDetail, 0, 100), // 保留最近100条记录
	}
}

// RecordCall 记录API调用
func (s *Stats) RecordCall(path, method, ip string, statusCode int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 增加总调用次数
	s.TotalCalls++

	// 增加HTTP方法调用次数
	s.MethodCalls[method]++

	// 增加API路径调用次数
	s.PathCalls[path]++

	// 增加IP调用次数
	s.IPCalls[ip]++

	// 记录当前时间
	now := time.Now()

	// 检查是否需要重置日统计
	if now.Day() != s.LastResetTime.Day() {
		s.DailyCalls = 1
		s.LastResetTime = now
	} else {
		s.DailyCalls++
	}

	// 记录最近调用详情
	detail := &CallDetail{
		Path:       path,
		Method:     method,
		IP:         ip,
		Timestamp:  now,
		StatusCode: statusCode,
	}

	// 保持最多100条记录
	if len(s.LastCallDetails) >= 100 {
		// 移除最旧的记录
		s.LastCallDetails = s.LastCallDetails[1:]
	}
	s.LastCallDetails = append(s.LastCallDetails, detail)
}

// GetStats 获取统计信息
func (s *Stats) GetStats() *Stats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 创建一个副本返回，避免并发问题
	copy := &Stats{
		TotalCalls:      s.TotalCalls,
		DailyCalls:      s.DailyCalls,
		HourlyCalls:     s.HourlyCalls,
		MethodCalls:     make(map[string]int64),
		PathCalls:       make(map[string]int64),
		IPCalls:         make(map[string]int64),
		LastResetTime:   s.LastResetTime,
		LastCallDetails: make([]*CallDetail, len(s.LastCallDetails)),
	}

	// 复制map数据
	for k, v := range s.MethodCalls {
		copy.MethodCalls[k] = v
	}
	for k, v := range s.PathCalls {
		copy.PathCalls[k] = v
	}
	for k, v := range s.IPCalls {
		copy.IPCalls[k] = v
	}

	// 复制调用详情
	for i, detail := range s.LastCallDetails {
		copy.LastCallDetails[i] = &CallDetail{
			Path:       detail.Path,
			Method:     detail.Method,
			IP:         detail.IP,
			Timestamp:  detail.Timestamp,
			StatusCode: detail.StatusCode,
		}
	}

	return copy
}
