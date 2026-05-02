package monitor

import (
	"context"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
	"github.com/sirupsen/logrus"
)

// 监控相关常量
const (
	DefaultCollectInterval = 5 * time.Second  // 默认采集间隔
	DefaultMaxHistory      = 2880             // 默认最大历史记录数（4小时）
	DefaultHistoryLimit    = 60               // 默认返回历史数据条数
)

// SystemMetrics 系统性能指标
type SystemMetrics struct {
	Timestamp   int64       `json:"timestamp"`
	CPU         CPUMetrics  `json:"cpu"`
	Memory      MemMetrics  `json:"memory"`
	Network     NetMetrics  `json:"network"`
}

// CPUMetrics CPU指标
type CPUMetrics struct {
	TotalPercent  float64   `json:"total_percent"`
	CorePercents  []float64 `json:"core_percents"`
	CoreCount     int       `json:"core_count"`
}

// MemMetrics 内存指标
type MemMetrics struct {
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Free        uint64  `json:"free"`
	UsedPercent float64 `json:"used_percent"`
}

// NetMetrics 网络指标
type NetMetrics struct {
	UploadSpeed   float64 `json:"upload_speed"`   // bytes/s
	DownloadSpeed float64 `json:"download_speed"` // bytes/s
	TotalSent     uint64  `json:"total_sent"`
	TotalRecv     uint64  `json:"total_recv"`
}

// Monitor 系统监控器
type Monitor struct {
	mu            sync.RWMutex
	current       *SystemMetrics
	history       []SystemMetrics
	maxHistory    int
	lastNetIO     *net.IOCountersStat
	lastNetTime   time.Time
	stopCh        chan struct{}
	wg            sync.WaitGroup
}

// MonitorOptions 监控器配置选项
type MonitorOptions struct {
	CollectInterval time.Duration
	MaxHistory      int
}

// NewMonitor 创建监控器
func NewMonitor(opts ...MonitorOptions) *Monitor {
	opt := MonitorOptions{
		CollectInterval: DefaultCollectInterval,
		MaxHistory:      DefaultMaxHistory,
	}
	if len(opts) > 0 {
		opt = opts[0]
	}

	return &Monitor{
		current:     &SystemMetrics{},
		history:     make([]SystemMetrics, 0, opt.MaxHistory),
		maxHistory:  opt.MaxHistory,
		stopCh:      make(chan struct{}),
	}
}

// Start 开始监控
func (m *Monitor) Start() {
	m.collectWithTimeout()
	
	ticker := time.NewTicker(DefaultCollectInterval)
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		for {
			select {
			case <-ticker.C:
				m.collectWithTimeout()
			case <-m.stopCh:
				ticker.Stop()
				return
			}
		}
	}()
}

// Stop 停止监控
func (m *Monitor) Stop() {
	close(m.stopCh)
	m.wg.Wait()
}

// GetCurrent 获取当前指标（返回副本，避免外部修改）
func (m *Monitor) GetCurrent() SystemMetrics {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.current == nil {
		return SystemMetrics{}
	}
	
	// 深拷贝，避免外部修改影响内部数据
	result := *m.current
	if len(m.current.CPU.CorePercents) > 0 {
		result.CPU.CorePercents = make([]float64, len(m.current.CPU.CorePercents))
		copy(result.CPU.CorePercents, m.current.CPU.CorePercents)
	}
	return result
}

// GetHistory 获取历史数据（返回副本）
func (m *Monitor) GetHistory(limit int) []SystemMetrics {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if limit <= 0 {
		limit = DefaultHistoryLimit
	}
	if limit > len(m.history) {
		limit = len(m.history)
	}

	result := make([]SystemMetrics, limit)
	startIdx := len(m.history) - limit
	for i := 0; i < limit; i++ {
		result[i] = m.history[startIdx+i]
		// 深拷贝核心数据
		if len(m.history[startIdx+i].CPU.CorePercents) > 0 {
			result[i].CPU.CorePercents = make([]float64, len(m.history[startIdx+i].CPU.CorePercents))
			copy(result[i].CPU.CorePercents, m.history[startIdx+i].CPU.CorePercents)
		}
	}
	return result
}

// collectWithTimeout 带超时的数据采集
func (m *Monitor) collectWithTimeout() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	done := make(chan struct{})
	go func() {
		m.collect()
		close(done)
	}()

	select {
	case <-done:
		// 采集成功
	case <-ctx.Done():
		logrus.Warn("系统指标采集超时")
	}
}

// collect 采集数据
func (m *Monitor) collect() {
	metrics := &SystemMetrics{
		Timestamp: time.Now().Unix(),
	}

	// CPU - 使用更短的采样时间减少阻塞
	if percents, err := cpu.Percent(100*time.Millisecond, true); err == nil {
		metrics.CPU.CorePercents = make([]float64, len(percents))
		copy(metrics.CPU.CorePercents, percents)
		metrics.CPU.CoreCount = len(percents)
		if len(percents) > 0 {
			var total float64
			for _, p := range percents {
				total += p
			}
			metrics.CPU.TotalPercent = total / float64(len(percents))
		}
	} else {
		logrus.Debugf("CPU采集失败: %v", err)
	}

	// Memory
	if vmStat, err := mem.VirtualMemory(); err == nil {
		metrics.Memory.Total = vmStat.Total
		metrics.Memory.Used = vmStat.Used
		metrics.Memory.Free = vmStat.Free
		metrics.Memory.UsedPercent = vmStat.UsedPercent
	} else {
		logrus.Debugf("内存采集失败: %v", err)
	}

	// Network
	if netIOs, err := net.IOCounters(false); err == nil && len(netIOs) > 0 {
		now := time.Now()
		currentIO := &netIOs[0]

		if m.lastNetIO != nil {
			duration := now.Sub(m.lastNetTime).Seconds()
			if duration > 0 {
				sentDiff := int64(currentIO.BytesSent) - int64(m.lastNetIO.BytesSent)
				recvDiff := int64(currentIO.BytesRecv) - int64(m.lastNetIO.BytesRecv)
				
				if sentDiff >= 0 {
					metrics.Network.UploadSpeed = float64(sentDiff) / duration
				}
				if recvDiff >= 0 {
					metrics.Network.DownloadSpeed = float64(recvDiff) / duration
				}
			}
		}

		metrics.Network.TotalSent = currentIO.BytesSent
		metrics.Network.TotalRecv = currentIO.BytesRecv

		m.lastNetIO = currentIO
		m.lastNetTime = now
	} else if err != nil {
		logrus.Debugf("网络采集失败: %v", err)
	}

	m.mu.Lock()
	m.current = metrics
	m.history = append(m.history, *metrics)
	if len(m.history) > m.maxHistory {
		// 使用切片复用减少内存分配
		copy(m.history, m.history[len(m.history)-m.maxHistory:])
		m.history = m.history[:m.maxHistory]
	}
	m.mu.Unlock()
}

// GlobalMonitor 全局监控实例
var GlobalMonitor = NewMonitor()

func init() {
	GlobalMonitor.Start()
}
