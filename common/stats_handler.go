package common

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 自定义模板函数
var templateFuncs = template.FuncMap{
	"percentage": func(total, count int64) int {
		if total == 0 {
			return 0
		}
		return int((float64(count) / float64(total)) * 100)
	},
}

// StatsHandler 处理统计信息展示页面
func StatsHandler(c *gin.Context) {
	// 获取统计信息
	stats := GlobalStats.GetStats()

	// 向上下文添加自定义模板函数
	c.HTML(http.StatusOK, "stats.html", gin.H{
		"Stats": stats,
		"Funcs": templateFuncs,
	})
}
