package models

import (
	"time"
)

// APIKey 表示API密钥的模型
type APIKey struct {
	ID           int64     `json:"id"`
	Key          string    `json:"key"`
	Name         string    `json:"name"`
	MaxUsage     int64     `json:"max_usage"`
	CurrentUsage int64     `json:"current_usage"`
	IsPermanent  bool      `json:"is_permanent"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// APIKeyResponse 表示API密钥响应的模型
type APIKeyResponse struct {
	Key          string `json:"key"`
	Name         string `json:"name"`
	MaxUsage     int64  `json:"max_usage"`
	CurrentUsage int64  `json:"current_usage"`
	IsPermanent  bool   `json:"is_permanent"`
}
