package common

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xrcuo/xrcuo-api/config"
)

// Response 统一响应结构体
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
	Took string      `json:"took,omitempty"`
}

// SuccessResponse 成功响应
func SuccessResponse(c *gin.Context, data interface{}, msg string) {
	response := &Response{
		Code: 200,
		Msg:  msg,
		Data: data,
	}
	c.JSON(http.StatusOK, response)
}

// ErrorResponse 错误响应
func ErrorResponse(c *gin.Context, statusCode int, code int, msg string) {
	response := &Response{
		Code: code,
		Msg:  msg,
	}
	c.JSON(statusCode, response)
}

// JSONResponse 根据配置返回格式化或非格式化的JSON响应
func JSONResponse(c *gin.Context, statusCode int, obj interface{}) {
	c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	c.Writer.WriteHeader(statusCode)
	
	encoder := json.NewEncoder(c.Writer)
	
	if config.IsJSONFormatEnabled() {
		// 如果启用了格式化，使用固定的两个空格缩进
		encoder.SetIndent("", "  ")
	}
	
	// 编码并写入响应
	encoder.Encode(obj)
}
