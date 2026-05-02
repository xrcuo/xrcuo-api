package handler

import (
	"strconv"

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
	limit := monitor.DefaultHistoryLimit
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}
	
	history := monitor.GlobalMonitor.GetHistory(limit)
	common.SuccessResponse(c, history, "获取成功")
}
