package random

import (
	"github.com/gin-gonic/gin"
)

// RegisterRouter 注册随机图片插件的路由
func RegisterRouter(router *gin.RouterGroup) {
	{
		// 获取随机图片（重定向到图片URL）
		router.GET("/random/image", GetRandomImageHandler)
		// 获取随机图片信息（返回JSON格式）
		router.GET("/random/image/info", GetRandomImageInfoHandler)
	}
}
