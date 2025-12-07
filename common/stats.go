package common

import (
	"log"
	"sync"
	"time"

	"github.com/xrcuo/xrcuo-api/db"
	"github.com/xrcuo/xrcuo-api/models"
)

// Stats 存储API调用统计信息
type Stats struct {
	models.Stats
	mu sync.RWMutex // 读写锁，保证并发安全
}

// 全局统计实例
var GlobalStats *Stats

// InitStats 初始化统计信息
func InitStats() {
	// 从数据库加载统计数据
	statsData, err := db.LoadStats()
	if err != nil {
		log.Printf("从数据库加载统计数据失败: %v，使用默认值", err)
		// 使用默认值
		statsData = &models.Stats{
			MethodCalls:     make(map[string]int64),
			PathCalls:       make(map[string]int64),
			IPCalls:         make(map[string]int64),
			LastResetTime:   time.Now(),
			LastCallDetails: make([]*models.CallDetail, 0, 100), // 保留最近100条记录
		}
	}

	stats := &Stats{
		Stats: *statsData,
	}

	GlobalStats = stats

	// 启动定时保存任务（每30秒保存一次统计数据）
	go startPeriodicSave()
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
	detail := &models.CallDetail{
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

	// 异步保存调用详情到数据库
	go func() {
		if err := db.SaveCallDetail(detail); err != nil {
			log.Printf("保存调用详情到数据库失败: %v", err)
		}
	}()
}

// SaveStats 保存统计信息到数据库
func (s *Stats) SaveStats() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 创建副本
	statsCopy := s.GetStats()

	// 保存到数据库
	return db.SaveStats(statsCopy)
}

// startPeriodicSave 启动定时保存任务
func startPeriodicSave() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if err := GlobalStats.SaveStats(); err != nil {
			log.Printf("定时保存统计数据失败: %v", err)
		} else {
			log.Println("统计数据已保存到数据库")
		}
	}
}

// GetStats 获取统计信息
func (s *Stats) GetStats() *models.Stats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 创建一个副本返回，避免并发问题
	copy := &models.Stats{
		TotalCalls:      s.TotalCalls,
		DailyCalls:      s.DailyCalls,
		HourlyCalls:     s.HourlyCalls,
		MethodCalls:     make(map[string]int64),
		PathCalls:       make(map[string]int64),
		IPCalls:         make(map[string]int64),
		LastResetTime:   s.LastResetTime,
		LastCallDetails: make([]*models.CallDetail, len(s.LastCallDetails)),
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
		copy.LastCallDetails[i] = &models.CallDetail{
			Path:       detail.Path,
			Method:     detail.Method,
			IP:         detail.IP,
			Timestamp:  detail.Timestamp,
			StatusCode: detail.StatusCode,
		}
	}

	return copy
}
