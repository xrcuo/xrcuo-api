package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/xrcuo/xrcuo-api/config"
	_ "modernc.org/sqlite"
)

var DB *sql.DB

// InitDB 初始化数据库连接
func InitDB() error {
	var err error

	// 获取配置的数据库路径
	dbPath := config.GetDatabasePath()

	// 创建或打开SQLite数据库文件
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("打开数据库失败: %v", err)
	}

	// 配置连接池
	DB.SetMaxOpenConns(config.GetMaxOpenConns()) // 最大打开连接数
	DB.SetMaxIdleConns(config.GetMaxIdleConns()) // 最大空闲连接数
	DB.SetConnMaxLifetime(-1)                    // 连接最大生命周期（-1表示无限制）
	DB.SetConnMaxIdleTime(10 * time.Minute)      // 空闲连接最大生命周期

	logrus.Debug("数据库连接池配置完成")

	// 测试数据库连接
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("连接数据库失败: %v", err)
	}

	// 创建表结构
	if err = createTables(); err != nil {
		return fmt.Errorf("创建表结构失败: %v", err)
	}

	logrus.Info("数据库初始化成功")
	return nil
}

// createTables 创建数据库表结构
func createTables() error {
	// 创建表结构的SQL语句列表
	createTableSQLs := []string{
		// 统计信息表
		`
		CREATE TABLE IF NOT EXISTS stats (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			total_calls INTEGER NOT NULL DEFAULT 0,
			daily_calls INTEGER NOT NULL DEFAULT 0,
			last_reset_time DATETIME NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
		`,
		// HTTP方法统计表
		`
		CREATE TABLE IF NOT EXISTS method_calls (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			method TEXT NOT NULL,
			count INTEGER NOT NULL DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(method)
		);
		`,
		// API路径统计表
		`
		CREATE TABLE IF NOT EXISTS path_calls (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			path TEXT NOT NULL,
			count INTEGER NOT NULL DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(path)
		);
		`,
		// IP调用统计表
		`
		CREATE TABLE IF NOT EXISTS ip_calls (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			ip TEXT NOT NULL,
			count INTEGER NOT NULL DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(ip)
		);
		`,
		// API调用详情表
		`
		CREATE TABLE IF NOT EXISTS call_details (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			path TEXT NOT NULL,
			method TEXT NOT NULL,
			ip TEXT NOT NULL,
			timestamp DATETIME NOT NULL,
			status_code INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
		`,
		// API密钥表
		`
		CREATE TABLE IF NOT EXISTS api_keys (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			key TEXT NOT NULL UNIQUE,
			name TEXT NOT NULL,
			max_usage INTEGER NOT NULL DEFAULT 0,
			current_usage INTEGER NOT NULL DEFAULT 0,
			is_permanent BOOLEAN NOT NULL DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
		`,
	}

	// 执行创建表结构的SQL语句
	for _, sql := range createTableSQLs {
		if _, err := DB.Exec(sql); err != nil {
			return fmt.Errorf("创建表结构失败: %v", err)
		}
	}

	// 创建索引以提高查询性能
	indexSQLs := []string{
		"CREATE INDEX IF NOT EXISTS idx_call_details_timestamp ON call_details(timestamp DESC);",
		"CREATE INDEX IF NOT EXISTS idx_call_details_path ON call_details(path);",
		"CREATE INDEX IF NOT EXISTS idx_call_details_method ON call_details(method);",
		"CREATE INDEX IF NOT EXISTS idx_call_details_status ON call_details(status_code);",
		"CREATE INDEX IF NOT EXISTS idx_api_keys_key ON api_keys(key);",
		"CREATE INDEX IF NOT EXISTS idx_ip_calls_ip ON ip_calls(ip);",
		"CREATE INDEX IF NOT EXISTS idx_path_calls_path ON path_calls(path);",
		"CREATE INDEX IF NOT EXISTS idx_method_calls_method ON method_calls(method);",
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
		logrus.Info("正在关闭数据库连接")
		return DB.Close()
	}
	return nil
}

// GetDB 获取数据库连接实例
func GetDB() *sql.DB {
	return DB
}

// Transaction 执行事务
func Transaction(fn func(tx *sql.Tx) error) error {
	// 开启事务
	tx, err := DB.Begin()
	if err != nil {
		return fmt.Errorf("开启事务失败: %v", err)
	}

	// 执行事务函数
	if err := fn(tx); err != nil {
		// 回滚事务
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("事务回滚失败: %v, 原始错误: %v", rollbackErr, err)
		}
		return err
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("提交事务失败: %v", err)
	}

	return nil
}

// WithTransaction 在事务中执行查询
func WithTransaction(tx *sql.Tx, query string, args ...interface{}) (*sql.Rows, error) {
	if tx != nil {
		return tx.Query(query, args...)
	}
	return DB.Query(query, args...)
}

// WithTransactionExec 在事务中执行修改操作
func WithTransactionExec(tx *sql.Tx, query string, args ...interface{}) (sql.Result, error) {
	if tx != nil {
		return tx.Exec(query, args...)
	}
	return DB.Exec(query, args...)
}
