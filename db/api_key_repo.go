package db

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"time"

	"github.com/xrcuo/xrcuo-api/models"
)

// generateAPIKey 生成随机API密钥
// 使用32字节的随机数据，转换为64位的十六进制字符串
func generateAPIKey() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}

// CreateAPIKey 创建一个新的API密钥
// name: 密钥名称
// maxUsage: 最大使用次数，0表示无限制
// isPermanent: 是否为永久密钥
func CreateAPIKey(name string, maxUsage int64, isPermanent bool) (*models.APIKey, error) {
	// 生成API密钥
	key, err := generateAPIKey()
	if err != nil {
		return nil, fmt.Errorf("生成API密钥失败: %v", err)
	}

	now := time.Now()
	// 插入到数据库
	result, err := DB.Exec(
		"INSERT INTO api_keys (key, name, max_usage, current_usage, is_permanent, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		key, name, maxUsage, 0, isPermanent, now, now,
	)
	if err != nil {
		return nil, fmt.Errorf("创建API密钥失败: %v", err)
	}

	// 获取插入的ID
	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("获取API密钥ID失败: %v", err)
	}

	// 返回新创建的API密钥
	return &models.APIKey{
		ID:           id,
		Key:          key,
		Name:         name,
		MaxUsage:     maxUsage,
		CurrentUsage: 0,
		IsPermanent:  isPermanent,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// GetAPIKeyByKey 通过密钥字符串获取API密钥信息
// key: API密钥字符串
func GetAPIKeyByKey(key string) (*models.APIKey, error) {
	apiKey := &models.APIKey{}
	err := DB.QueryRow(
		"SELECT id, key, name, max_usage, current_usage, is_permanent, created_at, updated_at FROM api_keys WHERE key = ?",
		key,
	).Scan(
		&apiKey.ID, &apiKey.Key, &apiKey.Name, &apiKey.MaxUsage,
		&apiKey.CurrentUsage, &apiKey.IsPermanent, &apiKey.CreatedAt, &apiKey.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("API密钥不存在")
		}
		return nil, fmt.Errorf("查询API密钥失败: %v", err)
	}

	return apiKey, nil
}

// UpdateAPIKeyUsage 更新API密钥使用次数
// 使用一条UPDATE语句确保原子性，避免竞态条件
// key: API密钥字符串
func UpdateAPIKeyUsage(key string) error {
	// 使用一条UPDATE语句完成检查和更新，避免竞态条件
	result, err := DB.Exec(
		"UPDATE api_keys SET current_usage = current_usage + 1, updated_at = ? WHERE key = ? AND (is_permanent = 1 OR current_usage < max_usage)",
		time.Now(), key,
	)
	if err != nil {
		return fmt.Errorf("更新API密钥使用次数失败: %v", err)
	}

	// 检查是否有行被更新
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("获取影响行数失败: %v", err)
	}

	if rowsAffected == 0 {
		// 检查是因为密钥不存在还是达到了使用上限
		var count int
		if err := DB.QueryRow("SELECT COUNT(*) FROM api_keys WHERE key = ?", key).Scan(&count); err != nil {
			return fmt.Errorf("检查API密钥是否存在失败: %v", err)
		}

		if count == 0 {
			return fmt.Errorf("API密钥不存在")
		}
		return fmt.Errorf("API密钥已达到使用上限")
	}

	return nil
}

// DeleteAPIKey 删除API密钥
// id: API密钥ID
func DeleteAPIKey(id int64) error {
	_, err := DB.Exec("DELETE FROM api_keys WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("删除API密钥失败: %v", err)
	}
	return nil
}

// GetAllAPIKeys 获取所有API密钥
// 按创建时间倒序排列
func GetAllAPIKeys() ([]*models.APIKey, error) {
	rows, err := DB.Query(
		"SELECT id, key, name, max_usage, current_usage, is_permanent, created_at, updated_at FROM api_keys ORDER BY created_at DESC",
	)
	if err != nil {
		return nil, fmt.Errorf("查询所有API密钥失败: %v", err)
	}
	defer rows.Close()

	var apiKeys []*models.APIKey
	for rows.Next() {
		apiKey := &models.APIKey{}
		if err := rows.Scan(
			&apiKey.ID, &apiKey.Key, &apiKey.Name, &apiKey.MaxUsage,
			&apiKey.CurrentUsage, &apiKey.IsPermanent, &apiKey.CreatedAt, &apiKey.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("扫描API密钥失败: %v", err)
		}
		apiKeys = append(apiKeys, apiKey)
	}

	return apiKeys, nil
}
