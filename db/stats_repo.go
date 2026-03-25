package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/xrcuo/xrcuo-api/config"
	"github.com/xrcuo/xrcuo-api/models"
)

func LoadStats() (*models.Stats, error) {
	stats := &models.Stats{
		MethodCalls: make(map[string]int64),
		PathCalls:   make(map[string]int64),
		IPCalls:     make(map[string]int64),
	}

	dbType := config.GetDatabaseType()
	var limitSQL string
	if dbType == "postgresql" {
		limitSQL = "LIMIT 1"
	} else {
		limitSQL = "LIMIT 1"
	}

	row := DB.QueryRow("SELECT total_calls, daily_calls, last_reset_time FROM stats ORDER BY updated_at DESC " + limitSQL)
	err := row.Scan(&stats.TotalCalls, &stats.DailyCalls, &stats.LastResetTime)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("加载基本统计信息失败: %v", err)
	}

	if err == sql.ErrNoRows {
		stats.TotalCalls = 0
		stats.DailyCalls = 0
		stats.LastResetTime = time.Now()
	}

	rows, err := DB.Query("SELECT method, count FROM method_calls")
	if err != nil {
		return nil, fmt.Errorf("加载HTTP方法统计失败: %v", err)
	}

	for rows.Next() {
		var method string
		var count int64
		if scanErr := rows.Scan(&method, &count); scanErr != nil {
			rows.Close()
			return nil, fmt.Errorf("扫描HTTP方法统计失败: %v", scanErr)
		}
		stats.MethodCalls[method] = count
	}
	rows.Close()

	rows, err = DB.Query("SELECT path, count FROM path_calls")
	if err != nil {
		return nil, fmt.Errorf("加载API路径统计失败: %v", err)
	}

	for rows.Next() {
		var path string
		var count int64
		if scanErr := rows.Scan(&path, &count); scanErr != nil {
			rows.Close()
			return nil, fmt.Errorf("扫描API路径统计失败: %v", scanErr)
		}
		stats.PathCalls[path] = count
	}
	rows.Close()

	rows, err = DB.Query("SELECT ip, count FROM ip_calls")
	if err != nil {
		return nil, fmt.Errorf("加载IP统计失败: %v", err)
	}

	for rows.Next() {
		var ip string
		var count int64
		if scanErr := rows.Scan(&ip, &count); scanErr != nil {
			rows.Close()
			return nil, fmt.Errorf("扫描IP统计失败: %v", scanErr)
		}
		stats.IPCalls[ip] = count
	}
	rows.Close()

	rows, err = DB.Query(
		"SELECT path, method, ip, timestamp, status_code FROM call_details ORDER BY timestamp DESC LIMIT 100",
	)
	if err != nil {
		return nil, fmt.Errorf("加载调用详情失败: %v", err)
	}

	details := make([]*models.CallDetail, 0, 100)
	for rows.Next() {
		var detail models.CallDetail
		if scanErr := rows.Scan(&detail.Path, &detail.Method, &detail.IP, &detail.Timestamp, &detail.StatusCode); scanErr != nil {
			rows.Close()
			return nil, fmt.Errorf("扫描调用详情失败: %v", scanErr)
		}
		details = append(details, &detail)
	}
	rows.Close()

	for i, j := 0, len(details)-1; i < j; i, j = i+1, j-1 {
		details[i], details[j] = details[j], details[i]
	}

	stats.LastCallDetails = details

	return stats, nil
}

func SaveStats(stats *models.Stats) error {
	tx, err := DB.Begin()
	if err != nil {
		return fmt.Errorf("开启事务失败: %v", err)
	}
	defer func() {
		if tx != nil {
			tx.Rollback()
		}
	}()

	if err = saveStatsTx(tx, stats); err != nil {
		return err
	}

	return tx.Commit()
}

func saveStatsTx(tx *sql.Tx, stats *models.Stats) error {
	dbType := config.GetDatabaseType()
	now := time.Now()

	var id int64
	err := tx.QueryRow("SELECT id FROM stats ORDER BY updated_at DESC LIMIT 1").Scan(&id)
	if err == sql.ErrNoRows {
		_, err = tx.Exec(
			"INSERT INTO stats (total_calls, daily_calls, last_reset_time, created_at, updated_at) VALUES ("+
				GetPlaceholder(1)+", "+GetPlaceholder(2)+", "+GetPlaceholder(3)+", "+GetPlaceholder(4)+", "+GetPlaceholder(5)+")",
			stats.TotalCalls, stats.DailyCalls, stats.LastResetTime, now, now,
		)
	} else {
		_, err = tx.Exec(
			"UPDATE stats SET total_calls="+GetPlaceholder(1)+", daily_calls="+GetPlaceholder(2)+
				", last_reset_time="+GetPlaceholder(3)+", updated_at="+GetPlaceholder(4)+" WHERE id="+GetPlaceholder(5),
			stats.TotalCalls, stats.DailyCalls, stats.LastResetTime, now, id,
		)
	}
	if err != nil {
		return fmt.Errorf("保存基本统计信息失败: %v", err)
	}

	for method, count := range stats.MethodCalls {
		err = upsertMethodCall(tx, method, count, dbType)
		if err != nil {
			return fmt.Errorf("保存HTTP方法统计失败: %v", err)
		}
	}

	for path, count := range stats.PathCalls {
		err = upsertPathCall(tx, path, count, dbType)
		if err != nil {
			return fmt.Errorf("保存API路径统计失败: %v", err)
		}
	}

	for ip, count := range stats.IPCalls {
		err = upsertIPCall(tx, ip, count, dbType)
		if err != nil {
			return fmt.Errorf("保存IP统计失败: %v", err)
		}
	}

	return nil
}

func upsertMethodCall(tx *sql.Tx, method string, count int64, dbType string) error {
	now := time.Now()
	var exists bool
	err := tx.QueryRow("SELECT COUNT(1) FROM method_calls WHERE method = "+GetPlaceholder(1), method).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		_, err = tx.Exec("UPDATE method_calls SET count="+GetPlaceholder(1)+", updated_at="+GetPlaceholder(2)+" WHERE method="+GetPlaceholder(3),
			count, now, method)
	} else {
		_, err = tx.Exec("INSERT INTO method_calls (method, count, created_at, updated_at) VALUES ("+
			GetPlaceholder(1)+", "+GetPlaceholder(2)+", "+GetPlaceholder(3)+", "+GetPlaceholder(4)+")",
			method, count, now, now)
	}
	return err
}

func upsertPathCall(tx *sql.Tx, path string, count int64, dbType string) error {
	now := time.Now()
	var exists bool
	err := tx.QueryRow("SELECT COUNT(1) FROM path_calls WHERE path = "+GetPlaceholder(1), path).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		_, err = tx.Exec("UPDATE path_calls SET count="+GetPlaceholder(1)+", updated_at="+GetPlaceholder(2)+" WHERE path="+GetPlaceholder(3),
			count, now, path)
	} else {
		_, err = tx.Exec("INSERT INTO path_calls (path, count, created_at, updated_at) VALUES ("+
			GetPlaceholder(1)+", "+GetPlaceholder(2)+", "+GetPlaceholder(3)+", "+GetPlaceholder(4)+")",
			path, count, now, now)
	}
	return err
}

func upsertIPCall(tx *sql.Tx, ip string, count int64, dbType string) error {
	now := time.Now()
	var exists bool
	err := tx.QueryRow("SELECT COUNT(1) FROM ip_calls WHERE ip = "+GetPlaceholder(1), ip).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		_, err = tx.Exec("UPDATE ip_calls SET count="+GetPlaceholder(1)+", updated_at="+GetPlaceholder(2)+" WHERE ip="+GetPlaceholder(3),
			count, now, ip)
	} else {
		_, err = tx.Exec("INSERT INTO ip_calls (ip, count, created_at, updated_at) VALUES ("+
			GetPlaceholder(1)+", "+GetPlaceholder(2)+", "+GetPlaceholder(3)+", "+GetPlaceholder(4)+")",
			ip, count, now, now)
	}
	return err
}

func SaveCallDetailsBatch(details []*models.CallDetail) error {
	if len(details) == 0 {
		return nil
	}

	tx, err := DB.Begin()
	if err != nil {
		return fmt.Errorf("开启事务失败: %v", err)
	}
	defer func() {
		if tx != nil {
			tx.Rollback()
		}
	}()

	stmt, err := tx.Prepare("INSERT INTO call_details (path, method, ip, timestamp, status_code) VALUES (" +
		GetPlaceholder(1) + ", " + GetPlaceholder(2) + ", " + GetPlaceholder(3) + ", " + GetPlaceholder(4) + ", " + GetPlaceholder(5) + ")")
	if err != nil {
		return fmt.Errorf("准备插入语句失败: %v", err)
	}
	defer stmt.Close()

	for _, detail := range details {
		_, err = stmt.Exec(detail.Path, detail.Method, detail.IP, detail.Timestamp, detail.StatusCode)
		if err != nil {
			return fmt.Errorf("插入调用详情失败: %v", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("提交事务失败: %v", err)
	}

	go func() {
		dbType := config.GetDatabaseType()
		if dbType == "sqlite" {
			DB.Exec("DELETE FROM call_details WHERE id NOT IN (SELECT id FROM call_details ORDER BY timestamp DESC LIMIT 1000)")
		} else if dbType == "mysql" {
			DB.Exec("DELETE cd1 FROM call_details cd1 LEFT JOIN (SELECT id FROM call_details ORDER BY timestamp DESC LIMIT 1000) cd2 ON cd1.id = cd2.id WHERE cd2.id IS NULL")
		} else if dbType == "postgresql" {
			DB.Exec("DELETE FROM call_details WHERE id NOT IN (SELECT id FROM call_details ORDER BY timestamp DESC LIMIT 1000)")
		}
	}()

	return nil
}
