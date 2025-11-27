package common

import (
	"fmt"
	"strconv"
)

// StrToInt 字符串转整数（含默认值）
func StrToInt(str string, defaultValue int) int {
	if str == "" {
		return defaultValue
	}
	val, err := strconv.Atoi(str)
	if err != nil {
		return defaultValue
	}
	return val
}

// FormatDelay 格式化延迟显示（微秒→毫秒，更友好）
func FormatDelay(delay int64) string {
	if delay == 0 {
		return "超时"
	}
	// 转换为毫秒（保留2位小数）
	ms := float64(delay) / 1000.0
	return fmt.Sprintf("%.2fms", ms)
}
