package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/xrcuo/xrcuo-api/config"
	_ "modernc.org/sqlite"
)

var DB *sql.DB

// InitDB 初始化数据库连接
func InitDB() error {
	var err error
	dbType := config.GetDatabaseType()

	switch dbType {
	case "sqlite":
		err = initSQLite()
	case "mysql":
		err = initMySQL()
	case "postgresql":
		err = initPostgreSQL()
	default:
		return fmt.Errorf("不支持的数据库类型: %s", dbType)
	}

	if err != nil {
		return err
	}

	// 配置连接池
	DB.SetMaxOpenConns(config.GetMaxOpenConns())
	DB.SetMaxIdleConns(config.GetMaxIdleConns())
	DB.SetConnMaxLifetime(-1)
	DB.SetConnMaxIdleTime(10 * time.Minute)

	logrus.Debug("数据库连接池配置完成")

	// 测试数据库连接
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("连接数据库失败: %v", err)
	}

	// 创建表结构
	if err = createTables(); err != nil {
		return fmt.Errorf("创建表结构失败: %v", err)
	}

	logrus.Infof("数据库初始化成功 (类型: %s)", dbType)
	return nil
}

func initSQLite() error {
	dbPath := config.GetDatabasePath()
	var err error
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("打开SQLite数据库失败: %v", err)
	}
	return nil
}

func initMySQL() error {
	host := config.GetDatabaseHost()
	port := config.GetDatabasePort()
	user := config.GetDatabaseUser()
	password := config.GetDatabasePassword()
	dbname := config.GetDatabaseName()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbname)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("打开MySQL数据库失败: %v", err)
	}
	return nil
}

func initPostgreSQL() error {
	host := config.GetDatabaseHost()
	port := config.GetDatabasePort()
	user := config.GetDatabaseUser()
	password := config.GetDatabasePassword()
	dbname := config.GetDatabaseName()

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("打开PostgreSQL数据库失败: %v", err)
	}
	return nil
}

func getSQLForTable(tableName string) string {
	dbType := config.GetDatabaseType()

	switch tableName {
	case "stats":
		return getStatsTableSQL(dbType)
	case "method_calls":
		return getMethodCallsTableSQL(dbType)
	case "path_calls":
		return getPathCallsTableSQL(dbType)
	case "ip_calls":
		return getIPCallsTableSQL(dbType)
	case "call_details":
		return getCallDetailsTableSQL(dbType)
	default:
		return ""
	}
}

func getStatsTableSQL(dbType string) string {
	switch dbType {
	case "sqlite":
		return `
		CREATE TABLE IF NOT EXISTS stats (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			total_calls INTEGER NOT NULL DEFAULT 0,
			daily_calls INTEGER NOT NULL DEFAULT 0,
			last_reset_time DATETIME NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
		`
	case "mysql":
		return `
		CREATE TABLE IF NOT EXISTS stats (
			id INT AUTO_INCREMENT PRIMARY KEY,
			total_calls INT NOT NULL DEFAULT 0,
			daily_calls INT NOT NULL DEFAULT 0,
			last_reset_time DATETIME NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
		`
	case "postgresql":
		return `
		CREATE TABLE IF NOT EXISTS stats (
			id SERIAL PRIMARY KEY,
			total_calls INTEGER NOT NULL DEFAULT 0,
			daily_calls INTEGER NOT NULL DEFAULT 0,
			last_reset_time TIMESTAMP NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		`
	default:
		return ""
	}
}

func getMethodCallsTableSQL(dbType string) string {
	switch dbType {
	case "sqlite":
		return `
		CREATE TABLE IF NOT EXISTS method_calls (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			method TEXT NOT NULL,
			count INTEGER NOT NULL DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(method)
		);
		`
	case "mysql":
		return `
		CREATE TABLE IF NOT EXISTS method_calls (
			id INT AUTO_INCREMENT PRIMARY KEY,
			method VARCHAR(255) NOT NULL,
			count INT NOT NULL DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			UNIQUE KEY (method)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
		`
	case "postgresql":
		return `
		CREATE TABLE IF NOT EXISTS method_calls (
			id SERIAL PRIMARY KEY,
			method VARCHAR(255) NOT NULL,
			count INTEGER NOT NULL DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(method)
		);
		`
	default:
		return ""
	}
}

func getPathCallsTableSQL(dbType string) string {
	switch dbType {
	case "sqlite":
		return `
		CREATE TABLE IF NOT EXISTS path_calls (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			path TEXT NOT NULL,
			count INTEGER NOT NULL DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(path)
		);
		`
	case "mysql":
		return `
		CREATE TABLE IF NOT EXISTS path_calls (
			id INT AUTO_INCREMENT PRIMARY KEY,
			path VARCHAR(500) NOT NULL,
			count INT NOT NULL DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			UNIQUE KEY (path)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
		`
	case "postgresql":
		return `
		CREATE TABLE IF NOT EXISTS path_calls (
			id SERIAL PRIMARY KEY,
			path VARCHAR(500) NOT NULL,
			count INTEGER NOT NULL DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(path)
		);
		`
	default:
		return ""
	}
}

func getIPCallsTableSQL(dbType string) string {
	switch dbType {
	case "sqlite":
		return `
		CREATE TABLE IF NOT EXISTS ip_calls (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			ip TEXT NOT NULL,
			count INTEGER NOT NULL DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(ip)
		);
		`
	case "mysql":
		return `
		CREATE TABLE IF NOT EXISTS ip_calls (
			id INT AUTO_INCREMENT PRIMARY KEY,
			ip VARCHAR(45) NOT NULL,
			count INT NOT NULL DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			UNIQUE KEY (ip)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
		`
	case "postgresql":
		return `
		CREATE TABLE IF NOT EXISTS ip_calls (
			id SERIAL PRIMARY KEY,
			ip VARCHAR(45) NOT NULL,
			count INTEGER NOT NULL DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(ip)
		);
		`
	default:
		return ""
	}
}

func getCallDetailsTableSQL(dbType string) string {
	switch dbType {
	case "sqlite":
		return `
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
	case "mysql":
		return `
		CREATE TABLE IF NOT EXISTS call_details (
			id INT AUTO_INCREMENT PRIMARY KEY,
			path VARCHAR(500) NOT NULL,
			method VARCHAR(255) NOT NULL,
			ip VARCHAR(45) NOT NULL,
			timestamp DATETIME NOT NULL,
			status_code INT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
		`
	case "postgresql":
		return `
		CREATE TABLE IF NOT EXISTS call_details (
			id SERIAL PRIMARY KEY,
			path VARCHAR(500) NOT NULL,
			method VARCHAR(255) NOT NULL,
			ip VARCHAR(45) NOT NULL,
			timestamp TIMESTAMP NOT NULL,
			status_code INTEGER NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		`
	default:
		return ""
	}
}

func getIndexSQLs() []string {
	dbType := config.GetDatabaseType()

	switch dbType {
	case "sqlite":
		return []string{
			"CREATE INDEX IF NOT EXISTS idx_call_details_timestamp ON call_details(timestamp DESC);",
			"CREATE INDEX IF NOT EXISTS idx_call_details_path ON call_details(path);",
			"CREATE INDEX IF NOT EXISTS idx_call_details_method ON call_details(method);",
			"CREATE INDEX IF NOT EXISTS idx_call_details_status ON call_details(status_code);",
			"CREATE INDEX IF NOT EXISTS idx_ip_calls_ip ON ip_calls(ip);",
			"CREATE INDEX IF NOT EXISTS idx_path_calls_path ON path_calls(path);",
			"CREATE INDEX IF NOT EXISTS idx_method_calls_method ON method_calls(method);",
		}
	case "mysql":
		return []string{
			"CREATE INDEX IF NOT EXISTS idx_call_details_timestamp ON call_details(timestamp DESC);",
			"CREATE INDEX IF NOT EXISTS idx_call_details_path ON call_details(path);",
			"CREATE INDEX IF NOT EXISTS idx_call_details_method ON call_details(method);",
			"CREATE INDEX IF NOT EXISTS idx_call_details_status ON call_details(status_code);",
			"CREATE INDEX IF NOT EXISTS idx_ip_calls_ip ON ip_calls(ip);",
			"CREATE INDEX IF NOT EXISTS idx_path_calls_path ON path_calls(path);",
			"CREATE INDEX IF NOT EXISTS idx_method_calls_method ON method_calls(method);",
		}
	case "postgresql":
		return []string{
			"CREATE INDEX IF NOT EXISTS idx_call_details_timestamp ON call_details(timestamp DESC);",
			"CREATE INDEX IF NOT EXISTS idx_call_details_path ON call_details(path);",
			"CREATE INDEX IF NOT EXISTS idx_call_details_method ON call_details(method);",
			"CREATE INDEX IF NOT EXISTS idx_call_details_status ON call_details(status_code);",
			"CREATE INDEX IF NOT EXISTS idx_ip_calls_ip ON ip_calls(ip);",
			"CREATE INDEX IF NOT EXISTS idx_path_calls_path ON path_calls(path);",
			"CREATE INDEX IF NOT EXISTS idx_method_calls_method ON method_calls(method);",
		}
	default:
		return []string{}
	}
}

func createTables() error {
	tables := []string{"stats", "method_calls", "path_calls", "ip_calls", "call_details"}

	for _, table := range tables {
		sql := getSQLForTable(table)
		if sql == "" {
			continue
		}
		if _, err := DB.Exec(sql); err != nil {
			return fmt.Errorf("创建表 %s 失败: %v", table, err)
		}
	}

	dbType := config.GetDatabaseType()
	if dbType == "mysql" {
		createMySQLIndexes()
	} else {
		indexSQLs := getIndexSQLs()
		for _, sql := range indexSQLs {
			if _, err := DB.Exec(sql); err != nil {
				logrus.Warnf("创建索引失败: %v", err)
			}
		}
	}

	return nil
}

func createMySQLIndexes() {
	indexes := []struct {
		table  string
		name   string
		column string
	}{
		{"call_details", "idx_call_details_timestamp", "timestamp DESC"},
		{"call_details", "idx_call_details_path", "path"},
		{"call_details", "idx_call_details_method", "method"},
		{"call_details", "idx_call_details_status", "status_code"},
		{"ip_calls", "idx_ip_calls_ip", "ip"},
		{"path_calls", "idx_path_calls_path", "path"},
		{"method_calls", "idx_method_calls_method", "method"},
	}

	for _, idx := range indexes {
		var exists bool
		err := DB.QueryRow(
			"SELECT COUNT(1) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = ? AND index_name = ?",
			idx.table, idx.name,
		).Scan(&exists)
		if err != nil {
			logrus.Warnf("检查索引 %s 存在性失败: %v", idx.name, err)
			continue
		}
		if exists {
			continue
		}

		_, err = DB.Exec(fmt.Sprintf("CREATE INDEX %s ON %s(%s)", idx.name, idx.table, idx.column))
		if err != nil {
			logrus.Warnf("创建索引 %s 失败: %v", idx.name, err)
		}
	}
}

func CloseDB() error {
	if DB != nil {
		logrus.Info("正在关闭数据库连接")
		return DB.Close()
	}
	return nil
}

func GetDB() *sql.DB {
	return DB
}

func Transaction(fn func(tx *sql.Tx) error) error {
	tx, err := DB.Begin()
	if err != nil {
		return fmt.Errorf("开启事务失败: %v", err)
	}

	if err := fn(tx); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("事务回滚失败: %v, 原始错误: %v", rollbackErr, err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("提交事务失败: %v", err)
	}

	return nil
}

func WithTransaction(tx *sql.Tx, query string, args ...interface{}) (*sql.Rows, error) {
	if tx != nil {
		return tx.Query(query, args...)
	}
	return DB.Query(query, args...)
}

func GetPlaceholder(n int) string {
	dbType := config.GetDatabaseType()
	if dbType == "postgresql" {
		return fmt.Sprintf("$%d", n)
	}
	return "?"
}

func WithTransactionExec(tx *sql.Tx, query string, args ...interface{}) (sql.Result, error) {
	if tx != nil {
		return tx.Exec(query, args...)
	}
	return DB.Exec(query, args...)
}
