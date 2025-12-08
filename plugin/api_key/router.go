package api_key

import (
	"github.com/gin-gonic/gin"
)

// RegisterRouter 注册API密钥管理路由
func RegisterRouter(r *gin.RouterGroup) {
	apiKeyGroup := r.Group("/api_key")
	{
		// 获取所有API密钥
		apiKeyGroup.GET("", GetAPIKeysHandler)
		// 创建新的API密钥
		apiKeyGroup.POST("", CreateAPIKeyHandler)
		// 删除API密钥
		apiKeyGroup.DELETE("/:id", DeleteAPIKeyHandler)
	}
}
