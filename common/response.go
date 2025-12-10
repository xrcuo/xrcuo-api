package common

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xrcuo/xrcuo-api/config"
)

// 标准化错误码
const (
	// 成功
	CodeSuccess = 200

	// 客户端错误
	CodeBadRequest       = 400 // 请求参数错误
	CodeUnauthorized     = 401 // 未授权
	CodeForbidden        = 403 // 禁止访问
	CodeNotFound         = 404 // 资源不存在
	CodeMethodNotAllowed = 405 // 方法不允许
	CodeTooManyRequests  = 429 // 请求过于频繁

	// 服务器错误
	CodeInternalServerError = 500 // 服务器内部错误
	CodeDatabaseError       = 501 // 数据库错误
	CodeCacheError          = 502 // 缓存错误
	CodeThirdPartyError     = 503 // 第三方服务错误

	// 业务错误
	CodeAPIKeyError     = 1001 // API密钥错误
	CodeIPError         = 1002 // IP相关错误
	CodeValidationError = 1003 // 数据验证错误
)

// ErrorType 错误类型
type ErrorType string

const (
	ErrorTypeClient     ErrorType = "client"     // 客户端错误
	ErrorTypeServer     ErrorType = "server"     // 服务器错误
	ErrorTypeBusiness   ErrorType = "business"   // 业务错误
	ErrorTypeSystem     ErrorType = "system"     // 系统错误
	ErrorTypeThirdParty ErrorType = "thirdparty" // 第三方服务错误
)

// AppError 应用错误结构体
type AppError struct {
	Code       int       `json:"code"`             // 错误码
	Message    string    `json:"message"`          // 错误消息
	Type       ErrorType `json:"type"`             // 错误类型
	Detail     string    `json:"detail,omitempty"` // 错误详情
	HTTPStatus int       `json:"http_status"`      // HTTP状态码
	Stack      string    `json:"stack,omitempty"`  // 堆栈信息（仅在调试模式）
}

// Error 实现error接口
func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s (code: %d, http: %d)", e.Type, e.Message, e.Code, e.HTTPStatus)
}

// Response 统一响应结构体
type Response struct {
	Code  int         `json:"code"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data,omitempty"`
	Took  string      `json:"took,omitempty"`
	Error *ErrorInfo  `json:"error,omitempty"`
}

// ErrorInfo 错误信息结构体
type ErrorInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Type    string `json:"type"`
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

// NewAppError 创建应用错误
func NewAppError(code int, message string, errType ErrorType, httpStatus int, detail string) *AppError {
	err := &AppError{
		Code:       code,
		Message:    message,
		Type:       errType,
		Detail:     detail,
		HTTPStatus: httpStatus,
	}

	// 在调试模式下，添加堆栈信息
	if config.GetServerMode() == gin.DebugMode {
		err.Stack = string(debug.Stack())
	}

	return err
}

// NewClientError 创建客户端错误
func NewClientError(code int, message string, detail string) *AppError {
	return NewAppError(code, message, ErrorTypeClient, http.StatusBadRequest, detail)
}

// NewServerError 创建服务器错误
func NewServerError(code int, message string, detail string) *AppError {
	return NewAppError(code, message, ErrorTypeServer, http.StatusInternalServerError, detail)
}

// NewBusinessError 创建业务错误
func NewBusinessError(code int, message string, detail string) *AppError {
	return NewAppError(code, message, ErrorTypeBusiness, http.StatusOK, detail)
}

// HandleAppError 处理应用错误
func HandleAppError(c *gin.Context, err error) {
	// 记录错误日志
	logrus.WithError(err).Error("处理请求时发生错误")

	// 如果是AppError类型，直接使用
	if appErr, ok := err.(*AppError); ok {
		response := &Response{
			Code: appErr.Code,
			Msg:  appErr.Message,
			Error: &ErrorInfo{
				Code:    appErr.Code,
				Message: appErr.Message,
				Type:    string(appErr.Type),
			},
		}
		c.JSON(appErr.HTTPStatus, response)
		return
	}

	// 否则，创建默认的服务器错误
	defaultErr := NewServerError(CodeInternalServerError, "服务器内部错误", err.Error())
	response := &Response{
		Code: defaultErr.Code,
		Msg:  defaultErr.Message,
		Error: &ErrorInfo{
			Code:    defaultErr.Code,
			Message: defaultErr.Message,
			Type:    string(defaultErr.Type),
		},
	}
	c.JSON(defaultErr.HTTPStatus, response)
}

// RecoveryMiddleware 全局错误恢复中间件
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				// 记录堆栈信息
				stack := string(debug.Stack())
				logrus.WithFields(logrus.Fields{
					"recover": r,
					"stack":   stack,
					"path":    c.Request.URL.Path,
					"method":  c.Request.Method,
				}).Error("发生panic")

				// 创建panic错误
				panicErr := NewAppError(
					CodeInternalServerError,
					"服务器内部错误",
					ErrorTypeSystem,
					http.StatusInternalServerError,
					fmt.Sprintf("%v", r),
				)
				panicErr.Stack = stack

				// 返回错误响应
				response := &Response{
					Code: panicErr.Code,
					Msg:  panicErr.Message,
					Error: &ErrorInfo{
						Code:    panicErr.Code,
						Message: panicErr.Message,
						Type:    string(panicErr.Type),
					},
				}
				c.JSON(panicErr.HTTPStatus, response)
				c.Abort()
			}
		}()

		c.Next()
	}
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
