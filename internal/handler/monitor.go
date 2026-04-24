package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/xrcuo/xrcuo-api/internal/monitor"
	"github.com/xrcuo/xrcuo-lib/common"
)

// GetSystemMetrics 获取当前系统指标
func GetSystemMetrics(c *gin.Context) {
	metrics := monitor.GlobalMonitor.GetCurrent()
	common.SuccessResponse(c, metrics, "获取成功")
}

// GetSystemMetricsHistory 获取历史指标数据
func GetSystemMetricsHistory(c *gin.Context) {
	limit := 60 // 默认返回最近60个点（5分钟）
	history := monitor.GlobalMonitor.GetHistory(limit)
	common.SuccessResponse(c, history, "获取成功")
}
