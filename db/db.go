package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

// InitDB 初始化数据库连接
func InitDB() error {
	var err error
	
	// 创建或打开SQLite数据库文件
	DB, err = sql.Open("sqlite", "./stats.db")
	if err != nil {
		return fmt.Errorf("打开数据库失败: %v", err)
	}

	// 测试数据库连接
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("连接数据库失败: %v", err)
	}

	// 创建表结构
	if err = createTables(); err != nil {
		return fmt.Errorf("创建表结构失败: %v", err)
	}

	log.Println("数据库初始化成功")
	return nil
}

// createTables 创建数据库表结构
func createTables() error {
	// 创建统计信息表
	statsTableSQL := `
	CREATE TABLE IF NOT EXISTS stats (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		total_calls INTEGER NOT NULL DEFAULT 0,
		daily_calls INTEGER NOT NULL DEFAULT 0,
		last_reset_time DATETIME NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	if _, err := DB.Exec(statsTableSQL); err != nil {
		return fmt.Errorf("创建stats表失败: %v", err)
	}

	// 创建HTTP方法统计表
	methodCallsTableSQL := `
	CREATE TABLE IF NOT EXISTS method_calls (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		method TEXT NOT NULL,
		count INTEGER NOT NULL DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(method)
	);
	`

	if _, err := DB.Exec(methodCallsTableSQL); err != nil {
		return fmt.Errorf("创建method_calls表失败: %v", err)
	}

	// 创建API路径统计表
	pathCallsTableSQL := `
	CREATE TABLE IF NOT EXISTS path_calls (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		path TEXT NOT NULL,
		count INTEGER NOT NULL DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(path)
	);
	`

	if _, err := DB.Exec(pathCallsTableSQL); err != nil {
		return fmt.Errorf("创建path_calls表失败: %v", err)
	}

	// 创建IP调用统计表
	ipCallsTableSQL := `
	CREATE TABLE IF NOT EXISTS ip_calls (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		ip TEXT NOT NULL,
		count INTEGER NOT NULL DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(ip)
	);
	`

	if _, err := DB.Exec(ipCallsTableSQL); err != nil {
		return fmt.Errorf("创建ip_calls表失败: %v", err)
	}

	// 创建API调用详情表
	callDetailsTableSQL := `
	CREATE TABLE IF NOT EXISTS call_details (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		path TEXT NOT NULL,
		method TEXT NOT NULL,
		ip TEXT NOT NULL,
		timestamp DATETIME NOT NULL,
		status_code INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	if _, err := DB.Exec(callDetailsTableSQL); err != nil {
		return fmt.Errorf("创建call_details表失败: %v", err)
	}

	// 创建索引以提高查询性能
	indexSQLs := []string{
		"CREATE INDEX IF NOT EXISTS idx_call_details_timestamp ON call_details(timestamp DESC);",
		"CREATE INDEX IF NOT EXISTS idx_call_details_path ON call_details(path);",
		"CREATE INDEX IF NOT EXISTS idx_call_details_method ON call_details(method);",
	}

	for _, sql := range indexSQLs {
		if _, err := DB.Exec(sql); err != nil {
			return fmt.Errorf("创建索引失败: %v", err)
		}
	}

	return nil
}

// CloseDB 关闭数据库连接
func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
