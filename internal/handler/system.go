package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xrcuo/xrcuo-api/internal/models"
	"github.com/xrcuo/xrcuo-lib/common"
)

// SystemConfigSetRequest 设置配置请求
type SystemConfigSetRequest struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value"`
}

// SystemConfigGet 获取配置
func SystemConfigGet(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		common.ErrorResponse(c, http.StatusBadRequest, common.CodeBadRequest, "缺少key参数")
		return
	}

	value, err := models.SystemConfigGet(key)
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, common.CodeInternalServerError, "获取配置失败")
		return
	}

	common.SuccessResponse(c, gin.H{
		"key":   key,
		"value": value,
	}, "获取成功")
}

// SystemConfigSet 设置配置
func SystemConfigSet(c *gin.Context) {
	var req SystemConfigSetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, common.CodeBadRequest, "请求参数错误")
		return
	}

	if err := models.SystemConfigSet(req.Key, req.Value); err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, common.CodeInternalServerError, "设置配置失败")
		return
	}

	common.SuccessResponse(c, nil, "设置成功")
}

// SystemConfigList 获取所有配置
func SystemConfigList(c *gin.Context) {
	configs, err := models.SystemConfigList()
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, common.CodeInternalServerError, "获取配置列表失败")
		return
	}

	common.SuccessResponse(c, configs, "获取成功")
}
