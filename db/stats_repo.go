package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/xrcuo/xrcuo-api/models"
)

// LoadStats 从数据库加载统计信息
func LoadStats() (*models.Stats, error) {
	stats := &models.Stats{
		MethodCalls: make(map[string]int64),
		PathCalls:   make(map[string]int64),
		IPCalls:     make(map[string]int64),
	}

	// 加载基本统计信息
	row := DB.QueryRow("SELECT total_calls, daily_calls, last_reset_time FROM stats ORDER BY updated_at DESC LIMIT 1")
	err := row.Scan(&stats.TotalCalls, &stats.DailyCalls, &stats.LastResetTime)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("加载基本统计信息失败: %v", err)
	}

	// 如果没有统计数据，初始化默认值
	if err == sql.ErrNoRows {
		stats.TotalCalls = 0
		stats.DailyCalls = 0
		stats.LastResetTime = time.Now()
	}

	// 加载HTTP方法统计
	rows, err := DB.Query("SELECT method, count FROM method_calls")
	if err != nil {
		return nil, fmt.Errorf("加载HTTP方法统计失败: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var method string
		var count int64
		if scanErr := rows.Scan(&method, &count); scanErr != nil {
			return nil, fmt.Errorf("扫描HTTP方法统计失败: %v", scanErr)
		}
		stats.MethodCalls[method] = count
	}

	// 加载API路径统计
	rows, err = DB.Query("SELECT path, count FROM path_calls")
	if err != nil {
		return nil, fmt.Errorf("加载API路径统计失败: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var path string
		var count int64
		if scanErr := rows.Scan(&path, &count); scanErr != nil {
			return nil, fmt.Errorf("扫描API路径统计失败: %v", scanErr)
		}
		stats.PathCalls[path] = count
	}

	// 加载IP统计
	rows, err = DB.Query("SELECT ip, count FROM ip_calls")
	if err != nil {
		return nil, fmt.Errorf("加载IP统计失败: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var ip string
		var count int64
		if scanErr := rows.Scan(&ip, &count); scanErr != nil {
			return nil, fmt.Errorf("扫描IP统计失败: %v", scanErr)
		}
		stats.IPCalls[ip] = count
	}

	// 加载最近的调用详情（最多100条）
	rows, err = DB.Query(
		"SELECT path, method, ip, timestamp, status_code FROM call_details ORDER BY timestamp DESC LIMIT 100",
	)
	if err != nil {
		return nil, fmt.Errorf("加载调用详情失败: %v", err)
	}
	defer rows.Close()

	details := make([]*models.CallDetail, 0, 100)
	for rows.Next() {
		var detail models.CallDetail
		if scanErr := rows.Scan(&detail.Path, &detail.Method, &detail.IP, &detail.Timestamp, &detail.StatusCode); scanErr != nil {
			return nil, fmt.Errorf("扫描调用详情失败: %v", scanErr)
		}
		details = append(details, &detail)
	}

	// 反转顺序，使最新的记录在最后
	for i, j := 0, len(details)-1; i < j; i, j = i+1, j-1 {
		details[i], details[j] = details[j], details[i]
	}

	stats.LastCallDetails = details

	return stats, nil
}

// SaveStats 保存统计信息到数据库
func SaveStats(stats *models.Stats) error {
	// 开启事务
	tx, err := DB.Begin()
	if err != nil {
		return fmt.Errorf("开启事务失败: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 保存基本统计信息
	_, err = tx.Exec(
		"INSERT OR REPLACE INTO stats (id, total_calls, daily_calls, last_reset_time, updated_at) "+
			"VALUES ((SELECT id FROM stats ORDER BY updated_at DESC LIMIT 1), ?, ?, ?, ?)",
		stats.TotalCalls, stats.DailyCalls, stats.LastResetTime, time.Now(),
	)
	if err != nil {
		return fmt.Errorf("保存基本统计信息失败: %v", err)
	}

	// 保存HTTP方法统计
	for method, count := range stats.MethodCalls {
		_, err = tx.Exec(
			"INSERT OR REPLACE INTO method_calls (method, count, updated_at) VALUES (?, ?, ?)",
			method, count, time.Now(),
		)
		if err != nil {
			return fmt.Errorf("保存HTTP方法统计失败: %v", err)
		}
	}

	// 保存API路径统计
	for path, count := range stats.PathCalls {
		_, err = tx.Exec(
			"INSERT OR REPLACE INTO path_calls (path, count, updated_at) VALUES (?, ?, ?)",
			path, count, time.Now(),
		)
		if err != nil {
			return fmt.Errorf("保存API路径统计失败: %v", err)
		}
	}

	// 保存IP统计
	for ip, count := range stats.IPCalls {
		_, err = tx.Exec(
			"INSERT OR REPLACE INTO ip_calls (ip, count, updated_at) VALUES (?, ?, ?)",
			ip, count, time.Now(),
		)
		if err != nil {
			return fmt.Errorf("保存IP统计失败: %v", err)
		}
	}

	// 提交事务
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("提交事务失败: %v", err)
	}

	return nil
}

// SaveCallDetail 保存API调用详情到数据库
func SaveCallDetail(detail *models.CallDetail) error {
	return SaveCallDetailsBatch([]*models.CallDetail{detail})
}

// SaveCallDetailsBatch 批量保存API调用详情到数据库
func SaveCallDetailsBatch(details []*models.CallDetail) error {
	if len(details) == 0 {
		return nil
	}

	// 开启事务
	tx, err := DB.Begin()
	if err != nil {
		return fmt.Errorf("开启事务失败: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 准备插入语句
	stmt, err := tx.Prepare("INSERT INTO call_details (path, method, ip, timestamp, status_code) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("准备插入语句失败: %v", err)
	}
	defer stmt.Close()

	// 批量插入数据
	for _, detail := range details {
		_, err = stmt.Exec(detail.Path, detail.Method, detail.IP, detail.Timestamp, detail.StatusCode)
		if err != nil {
			return fmt.Errorf("插入调用详情失败: %v", err)
		}
	}

	// 提交事务
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("提交事务失败: %v", err)
	}

	// 保留最近1000条记录，删除旧记录（异步执行，不阻塞主流程）
	go func() {
		_, err := DB.Exec(
			"DELETE FROM call_details WHERE id NOT IN (SELECT id FROM call_details ORDER BY timestamp DESC LIMIT 1000)",
		)
		if err != nil {
			log.Printf("清理旧调用记录失败: %v", err)
			// 不影响主流程，继续执行
		}
	}()

	return nil
}
