package db

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/xrcuo/xrcuo-api/models"
)

// 生成随机API密钥
func generateAPIKey() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}

// CreateAPIKey 创建一个新的API密钥
func CreateAPIKey(name string, maxUsage int64, isPermanent bool) (*models.APIKey, error) {
	// 生成API密钥
	key, err := generateAPIKey()
	if err != nil {
		return nil, fmt.Errorf("生成API密钥失败: %v", err)
	}

	// 插入到数据库
	result, err := DB.Exec(
		"INSERT INTO api_keys (key, name, max_usage, current_usage, is_permanent, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		key, name, maxUsage, 0, isPermanent, time.Now(), time.Now(),
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
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}

// GetAPIKeyByKey 通过密钥获取API密钥信息
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
		return nil, fmt.Errorf("查询API密钥失败: %v", err)
	}

	return apiKey, nil
}

// UpdateAPIKeyUsage 更新API密钥使用次数
func UpdateAPIKeyUsage(key string) error {
	// 检查API密钥是否存在且有效
	apiKey, err := GetAPIKeyByKey(key)
	if err != nil {
		return fmt.Errorf("API密钥不存在: %v", err)
	}

	// 检查API密钥是否已达到使用上限
	if !apiKey.IsPermanent && apiKey.CurrentUsage >= apiKey.MaxUsage {
		return fmt.Errorf("API密钥已达到使用上限")
	}

	// 更新使用次数
	_, err = DB.Exec(
		"UPDATE api_keys SET current_usage = current_usage + 1, updated_at = ? WHERE key = ?",
		time.Now(), key,
	)
	if err != nil {
		return fmt.Errorf("更新API密钥使用次数失败: %v", err)
	}

	return nil
}

// DeleteAPIKey 删除API密钥
func DeleteAPIKey(id int64) error {
	_, err := DB.Exec(
		"DELETE FROM api_keys WHERE id = ?",
		id,
	)
	if err != nil {
		return fmt.Errorf("删除API密钥失败: %v", err)
	}

	return nil
}

// GetAllAPIKeys 获取所有API密钥
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
