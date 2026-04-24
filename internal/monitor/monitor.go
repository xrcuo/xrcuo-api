package monitor

import (
	"sync"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
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
	UploadSpeed   float64 `json:"upload_speed"`   // MB/s
	DownloadSpeed float64 `json:"download_speed"` // MB/s
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
}

// NewMonitor 创建监控器
func NewMonitor() *Monitor {
	return &Monitor{
		current:     &SystemMetrics{},
		history:     make([]SystemMetrics, 0, 2880),
		maxHistory:  2880,
		stopCh:      make(chan struct{}),
	}
}

// Start 开始监控
func (m *Monitor) Start() {
	m.collect()
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				m.collect()
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
}

// GetCurrent 获取当前指标
func (m *Monitor) GetCurrent() *SystemMetrics {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.current == nil {
		return &SystemMetrics{}
	}
	return m.current
}

// GetHistory 获取历史数据
func (m *Monitor) GetHistory(limit int) []SystemMetrics {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if limit <= 0 || limit > len(m.history) {
		limit = len(m.history)
	}

	result := make([]SystemMetrics, limit)
	copy(result, m.history[len(m.history)-limit:])
	return result
}

// collect 采集数据
func (m *Monitor) collect() {
	metrics := &SystemMetrics{
		Timestamp: time.Now().Unix(),
	}

	// CPU
	if percents, err := cpu.Percent(0, true); err == nil {
		metrics.CPU.CorePercents = percents
		metrics.CPU.CoreCount = len(percents)
		if len(percents) > 0 {
			var total float64
			for _, p := range percents {
				total += p
			}
			metrics.CPU.TotalPercent = total / float64(len(percents))
		}
	}

	// Memory
	if vmStat, err := mem.VirtualMemory(); err == nil {
		metrics.Memory.Total = vmStat.Total
		metrics.Memory.Used = vmStat.Used
		metrics.Memory.Free = vmStat.Free
		metrics.Memory.UsedPercent = vmStat.UsedPercent
	}

	// Network
	if netIOs, err := net.IOCounters(false); err == nil && len(netIOs) > 0 {
		now := time.Now()
		currentIO := &netIOs[0]

		if m.lastNetIO != nil {
			duration := now.Sub(m.lastNetTime).Seconds()
			if duration > 0 {
				metrics.Network.UploadSpeed = float64(currentIO.BytesSent-m.lastNetIO.BytesSent) / duration / 1024 / 1024
				metrics.Network.DownloadSpeed = float64(currentIO.BytesRecv-m.lastNetIO.BytesRecv) / duration / 1024 / 1024
			}
		}

		metrics.Network.TotalSent = currentIO.BytesSent
		metrics.Network.TotalRecv = currentIO.BytesRecv

		m.lastNetIO = currentIO
		m.lastNetTime = now
	}

	m.mu.Lock()
	m.current = metrics
	m.history = append(m.history, *metrics)
	if len(m.history) > m.maxHistory {
		m.history = m.history[len(m.history)-m.maxHistory:]
	}
	m.mu.Unlock()
}

// GlobalMonitor 全局监控实例
var GlobalMonitor = NewMonitor()

func init() {
	GlobalMonitor.Start()
}
