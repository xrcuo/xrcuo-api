package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xrcuo/xrcuo-api/internal/auth"
	"github.com/xrcuo/xrcuo-api/internal/models"
	"github.com/xrcuo/xrcuo-lib/common"
)

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// Login 用户登录
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, common.CodeBadRequest, "请求参数错误: "+err.Error())
		return
	}

	user, err := models.UserVerifyPassword(req.Username, req.Password)
	if err != nil {
		common.ErrorResponse(c, http.StatusUnauthorized, common.CodeUnauthorized, "用户名或密码错误")
		return
	}

	token, err := auth.GenerateToken(user.Username, 24)
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, common.CodeInternalServerError, "生成Token失败")
		return
	}

	common.SuccessResponse(c, LoginResponse{
		Token:    token,
		Username: user.Username,
	}, "登录成功")
}

// ChangePassword 修改密码
func ChangePassword(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		common.ErrorResponse(c, http.StatusUnauthorized, common.CodeUnauthorized, "未登录")
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, common.CodeBadRequest, "请求参数错误: "+err.Error())
		return
	}

	err := models.UserUpdatePassword(username.(string), req.OldPassword, req.NewPassword)
	if err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, common.CodeBadRequest, err.Error())
		return
	}

	common.SuccessResponse(c, nil, "密码修改成功")
}

// GetCurrentUser 获取当前用户信息
func GetCurrentUser(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		common.ErrorResponse(c, http.StatusUnauthorized, common.CodeUnauthorized, "未登录")
		return
	}

	common.SuccessResponse(c, gin.H{
		"username": username.(string),
	}, "获取成功")
}

// AuthMiddleware JWT认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			common.ErrorResponse(c, http.StatusUnauthorized, common.CodeUnauthorized, "缺少Authorization头")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			common.ErrorResponse(c, http.StatusUnauthorized, common.CodeUnauthorized, "Authorization格式错误")
			c.Abort()
			return
		}

		claims, err := auth.ParseToken(parts[1])
		if err != nil {
			common.ErrorResponse(c, http.StatusUnauthorized, common.CodeUnauthorized, "Token无效或已过期")
			c.Abort()
			return
		}

		c.Set("username", claims.Username)
		c.Next()
	}
}
