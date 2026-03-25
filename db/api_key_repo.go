package db

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"time"

	"github.com/xrcuo/xrcuo-api/config"
	"github.com/xrcuo/xrcuo-api/models"
)

func generateAPIKey() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}

func getQuotedKeyField() string {
	dbType := config.GetDatabaseType()
	if dbType == "mysql" {
		return "`key`"
	}
	return "\"key\""
}

func CreateAPIKey(name string, maxUsage int64, isPermanent bool) (*models.APIKey, error) {
	key, err := generateAPIKey()
	if err != nil {
		return nil, fmt.Errorf("生成API密钥失败: %v", err)
	}

	now := time.Now()
	keyField := getQuotedKeyField()

	result, err := DB.Exec(
		"INSERT INTO api_keys ("+keyField+", name, max_usage, current_usage, is_permanent, created_at, updated_at) VALUES ("+
			GetPlaceholder(1)+", "+GetPlaceholder(2)+", "+GetPlaceholder(3)+", "+GetPlaceholder(4)+", "+GetPlaceholder(5)+", "+GetPlaceholder(6)+", "+GetPlaceholder(7)+")",
		key, name, maxUsage, 0, isPermanent, now, now,
	)
	if err != nil {
		return nil, fmt.Errorf("创建API密钥失败: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		dbType := config.GetDatabaseType()
		if dbType == "postgresql" {
			id = 0
		} else {
			return nil, fmt.Errorf("获取API密钥ID失败: %v", err)
		}
	}

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

func GetAPIKeyByKey(key string) (*models.APIKey, error) {
	apiKey := &models.APIKey{}
	keyField := getQuotedKeyField()
	err := DB.QueryRow(
		"SELECT id, "+keyField+", name, max_usage, current_usage, is_permanent, created_at, updated_at FROM api_keys WHERE "+keyField+" = "+GetPlaceholder(1),
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

func UpdateAPIKeyUsage(key string) error {
	keyField := getQuotedKeyField()
	dbType := config.GetDatabaseType()

	var isPermanentVal interface{}
	if dbType == "sqlite" {
		isPermanentVal = 1
	} else {
		isPermanentVal = true
	}

	result, err := DB.Exec(
		"UPDATE api_keys SET current_usage = current_usage + 1, updated_at = "+GetPlaceholder(1)+
			" WHERE "+keyField+" = "+GetPlaceholder(2)+" AND (is_permanent = "+GetPlaceholder(3)+" OR current_usage < max_usage)",
		time.Now(), key, isPermanentVal,
	)
	if err != nil {
		return fmt.Errorf("更新API密钥使用次数失败: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("获取影响行数失败: %v", err)
	}

	if rowsAffected == 0 {
		var count int
		if err := DB.QueryRow("SELECT COUNT(*) FROM api_keys WHERE "+keyField+" = "+GetPlaceholder(1), key).Scan(&count); err != nil {
			return fmt.Errorf("检查API密钥是否存在失败: %v", err)
		}

		if count == 0 {
			return fmt.Errorf("API密钥不存在")
		}
		return fmt.Errorf("API密钥已达到使用上限")
	}

	return nil
}

func DeleteAPIKey(id int64) error {
	_, err := DB.Exec("DELETE FROM api_keys WHERE id = "+GetPlaceholder(1), id)
	if err != nil {
		return fmt.Errorf("删除API密钥失败: %v", err)
	}
	return nil
}

func GetAllAPIKeys() ([]*models.APIKey, error) {
	keyField := getQuotedKeyField()
	rows, err := DB.Query(
		"SELECT id, " + keyField + ", name, max_usage, current_usage, is_permanent, created_at, updated_at FROM api_keys ORDER BY created_at DESC",
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
