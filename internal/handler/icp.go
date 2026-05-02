package handler

import (
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/xrcuo/xrcuo-api/internal/models"
	"github.com/xrcuo/xrcuo-lib/common"
)

// ICP 相关常量
const (
	icpConfigKey     = "site_icp"
	icpMaxLength     = 100  // 备案号最大长度
	icpMinLength     = 5    // 备案号最小长度
)

// ICPResponse 备案号响应
type ICPResponse struct {
	Value string `json:"value"`
}

// ICPSetRequest 设置备案号请求
type ICPSetRequest struct {
	Value string `json:"value" binding:"required"`
}

// validateICP 验证备案号格式
func validateICP(value string) (string, bool) {
	value = strings.TrimSpace(value)
	
	if value == "" {
		return "", false
	}
	
	length := utf8.RuneCountInString(value)
	if length < icpMinLength {
		return "", false
	}
	if length > icpMaxLength {
		value = string([]rune(value)[:icpMaxLength])
	}
	
	return value, true
}

// ICPGet 获取备案号
func ICPGet(c *gin.Context) {
	value, err := models.SystemConfigGet(icpConfigKey)
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, common.CodeInternalServerError, "获取备案号失败，请稍后重试")
		return
	}

	common.SuccessResponse(c, ICPResponse{Value: value}, "获取成功")
}

// ICPSet 设置备案号
func ICPSet(c *gin.Context) {
	var req ICPSetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, common.CodeBadRequest, "请求参数错误，请检查输入")
		return
	}

	validatedValue, ok := validateICP(req.Value)
	if !ok {
		common.ErrorResponse(c, http.StatusBadRequest, common.CodeBadRequest, "备案号格式不正确，长度应在5-100个字符之间")
		return
	}

	if err := models.SystemConfigSet(icpConfigKey, validatedValue); err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, common.CodeInternalServerError, "保存备案号失败，请稍后重试")
		return
	}

	common.SuccessResponse(c, ICPResponse{Value: validatedValue}, "备案号保存成功")
}
