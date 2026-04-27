package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xrcuo/xrcuo-api/internal/models"
	"github.com/xrcuo/xrcuo-lib/common"
)

const icpConfigKey = "site_icp"

// ICPGetRequest 获取备案号请求
type ICPGetRequest struct {
}

// ICPSetRequest 设置备案号请求
type ICPSetRequest struct {
	Value string `json:"value" binding:"required"`
}

// ICPGet 获取备案号
func ICPGet(c *gin.Context) {
	value, err := models.SystemConfigGet(icpConfigKey)
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, common.CodeInternalServerError, "获取备案号失败")
		return
	}

	common.SuccessResponse(c, gin.H{
		"value": value,
	}, "获取成功")
}

// ICPSet 设置备案号
func ICPSet(c *gin.Context) {
	var req ICPSetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, common.CodeBadRequest, "请求参数错误")
		return
	}

	if err := models.SystemConfigSet(icpConfigKey, req.Value); err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, common.CodeInternalServerError, "设置备案号失败")
		return
	}

	common.SuccessResponse(c, nil, "设置成功")
}
