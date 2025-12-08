package api_key

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xrcuo/xrcuo-api/db"
)

// GetAPIKeysHandler 获取所有API密钥
func GetAPIKeysHandler(c *gin.Context) {
	// 获取所有API密钥
	apiKeys, err := db.GetAllAPIKeys()
	if err != nil {
		logrus.Errorf("获取API密钥列表失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取API密钥列表失败",
		})
		return
	}

	// 返回API密钥列表
	c.JSON(http.StatusOK, gin.H{
		"api_keys": apiKeys,
	})
}

// CreateAPIKeyHandler 创建新的API密钥
func CreateAPIKeyHandler(c *gin.Context) {
	// 从请求体中获取参数
	var req struct {
		Name         string `json:"name" binding:"required"`
		MaxUsage     int64  `json:"max_usage"`
		IsPermanent  bool   `json:"is_permanent"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数无效",
		})
		return
	}

	// 创建API密钥
	apiKey, err := db.CreateAPIKey(req.Name, req.MaxUsage, req.IsPermanent)
	if err != nil {
		logrus.Errorf("创建API密钥失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "创建API密钥失败",
		})
		return
	}

	// 返回新创建的API密钥
	c.JSON(http.StatusOK, gin.H{
		"api_key": apiKey,
	})
}

// DeleteAPIKeyHandler 删除API密钥
func DeleteAPIKeyHandler(c *gin.Context) {
	// 从URL参数中获取ID
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的API密钥ID",
		})
		return
	}

	// 删除API密钥
	if err := db.DeleteAPIKey(id); err != nil {
		logrus.Errorf("删除API密钥失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "删除API密钥失败",
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"message": "API密钥删除成功",
	})
}
