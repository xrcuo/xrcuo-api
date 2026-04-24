package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/xrcuo/xrcuo-lib/db"
)

// SystemConfig 系统配置模型
type SystemConfig struct {
	ID        int64     `json:"id"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateSystemConfigTable 创建系统配置表
func CreateSystemConfigTable() error {
	sqlStr := `
	CREATE TABLE IF NOT EXISTS system_configs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		key TEXT NOT NULL UNIQUE,
		value TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := db.GetDB().Exec(sqlStr)
	return err
}

// SystemConfigGet 获取配置值
func SystemConfigGet(key string) (string, error) {
	var value string
	err := db.GetDB().QueryRow("SELECT value FROM system_configs WHERE key = ?", key).Scan(&value)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return value, nil
}

// SystemConfigSet 设置配置值
func SystemConfigSet(key, value string) error {
	_, err := db.GetDB().Exec(
		"INSERT INTO system_configs (key, value) VALUES (?, ?) ON CONFLICT(key) DO UPDATE SET value = excluded.value, updated_at = CURRENT_TIMESTAMP",
		key, value,
	)
	return err
}

// SystemConfigGetJSON 获取JSON配置
func SystemConfigGetJSON(key string, v interface{}) error {
	value, err := SystemConfigGet(key)
	if err != nil {
		return err
	}
	if value == "" {
		return nil
	}
	return json.Unmarshal([]byte(value), v)
}

// SystemConfigSetJSON 设置JSON配置
func SystemConfigSetJSON(key string, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return SystemConfigSet(key, string(data))
}

// SystemConfigList 获取所有配置
func SystemConfigList() ([]SystemConfig, error) {
	rows, err := db.GetDB().Query("SELECT id, key, value, created_at, updated_at FROM system_configs ORDER BY key")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configs []SystemConfig
	for rows.Next() {
		var cfg SystemConfig
		var createdAt, updatedAt sql.NullTime
		err := rows.Scan(&cfg.ID, &cfg.Key, &cfg.Value, &createdAt, &updatedAt)
		if err != nil {
			continue
		}
		if createdAt.Valid {
			cfg.CreatedAt = createdAt.Time
		}
		if updatedAt.Valid {
			cfg.UpdatedAt = updatedAt.Time
		}
		configs = append(configs, cfg)
	}
	return configs, nil
}

// InitAllTables 初始化所有表
func InitAllTables() error {
	if err := CreateUsersTable(); err != nil {
		return fmt.Errorf("创建用户表失败: %v", err)
	}
	if err := CreateApiDocsTable(); err != nil {
		return fmt.Errorf("创建API文档表失败: %v", err)
	}
	if err := CreateSystemConfigTable(); err != nil {
		return fmt.Errorf("创建系统配置表失败: %v", err)
	}
	return nil
}
