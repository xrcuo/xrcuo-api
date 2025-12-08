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

	// 调用详情缓冲区，用于批量写入数据库
	callDetailBuffer []*models.CallDetail
	bufferMutex      sync.Mutex // 缓冲区互斥锁
	bufferSize       int        // 缓冲区大小
	maxBufferSize    int        // 最大缓冲区大小
}

// 全局统计实例
var GlobalStats *Stats

// InitStats 初始化统计信息
func InitStats() {
	log.Println("初始化统计信息...")

	// 从数据库加载统计数据
	statsData, err := db.LoadStats()
	if err != nil {
		log.Printf("从数据库加载统计数据失败: %v，使用默认值", err)
		// 使用默认值初始化统计数据
		statsData = &models.Stats{
			TotalCalls:      0,
			DailyCalls:      0,
			HourlyCalls:     0,
			MethodCalls:     make(map[string]int64),
			PathCalls:       make(map[string]int64),
			IPCalls:         make(map[string]int64),
			LastResetTime:   time.Now(),
			LastCallDetails: make([]*models.CallDetail, 0, 100), // 保留最近100条记录
		}
	}

	// 创建并初始化统计实例
	stats := &Stats{
		Stats:            *statsData,
		callDetailBuffer: make([]*models.CallDetail, 0, 100), // 初始化缓冲区，容量100
		maxBufferSize:    100,                                // 最大缓冲区大小
	}

	// 设置全局统计实例
	GlobalStats = stats

	log.Println("统计信息初始化完成")

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

	// 将调用详情添加到缓冲区
	s.bufferMutex.Lock()
	s.callDetailBuffer = append(s.callDetailBuffer, detail)
	bufferSize := len(s.callDetailBuffer)
	s.bufferMutex.Unlock()

	// 当缓冲区达到最大大小时，异步批量写入数据库
	if bufferSize >= s.maxBufferSize {
		go s.flushCallDetailBuffer()
	}
}

// flushCallDetailBuffer 将缓冲区中的调用详情批量写入数据库
func (s *Stats) flushCallDetailBuffer() {
	s.bufferMutex.Lock()

	// 如果缓冲区为空，直接返回
	if len(s.callDetailBuffer) == 0 {
		s.bufferMutex.Unlock()
		return
	}

	// 将缓冲区中的数据复制到临时变量
	details := make([]*models.CallDetail, len(s.callDetailBuffer))
	copy(details, s.callDetailBuffer)

	// 清空缓冲区
	s.callDetailBuffer = s.callDetailBuffer[:0]
	s.bufferMutex.Unlock()

	// 批量写入数据库
	if err := db.SaveCallDetailsBatch(details); err != nil {
		log.Printf("批量保存调用详情到数据库失败: %v", err)
	}
}

// SaveStats 保存统计信息到数据库
func (s *Stats) SaveStats() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 创建副本
	statsCopy := s.GetStats()

	// 保存到数据库
	err := db.SaveStats(statsCopy)
	if err != nil {
		return err
	}

	// 同时将缓冲区中的调用详情写入数据库
	s.flushCallDetailBuffer()

	return nil
}

// startPeriodicSave 启动定时保存任务
func startPeriodicSave() {
	// 使用配置的保存间隔，默认30秒
	interval := 30 * time.Second

	log.Printf("启动统计数据定时保存任务，保存间隔: %v", interval)

	ticker := time.NewTicker(interval)
	defer func() {
		ticker.Stop()
		log.Println("统计数据定时保存任务已停止")
	}()

	for range ticker.C {
		// 检查 GlobalStats 是否为 nil
		if GlobalStats == nil {
			log.Println("统计数据实例未初始化，跳过保存")
			continue
		}

		// 执行保存操作
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
