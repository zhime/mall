package utils

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构体
type Response struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp string      `json:"timestamp"`
	RequestID string      `json:"request_id,omitempty"`
}

// 响应状态码常量
const (
	SUCCESS                = 200
	ERROR                  = 500
	INVALID_PARAMS         = 400
	UNAUTHORIZED           = 401
	FORBIDDEN              = 403
	NOT_FOUND              = 404
	METHOD_NOT_ALLOWED     = 405
	TOO_MANY_REQUESTS      = 429
	INTERNAL_SERVER_ERROR  = 500
	SERVICE_UNAVAILABLE    = 503
)

// 错误消息常量
const (
	MSG_SUCCESS                = "success"
	MSG_ERROR                  = "error"
	MSG_INVALID_PARAMS         = "invalid parameters"
	MSG_UNAUTHORIZED           = "unauthorized"
	MSG_FORBIDDEN              = "forbidden"
	MSG_NOT_FOUND              = "not found"
	MSG_METHOD_NOT_ALLOWED     = "method not allowed"
	MSG_TOO_MANY_REQUESTS      = "too many requests"
	MSG_INTERNAL_SERVER_ERROR  = "internal server error"
	MSG_SERVICE_UNAVAILABLE    = "service unavailable"
)

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	response := Response{
		Code:      SUCCESS,
		Message:   MSG_SUCCESS,
		Data:      data,
		Timestamp: time.Now().Format(time.RFC3339),
		RequestID: c.GetString("request_id"),
	}
	c.JSON(http.StatusOK, response)
}

// SuccessWithMessage 带消息的成功响应
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	response := Response{
		Code:      SUCCESS,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().Format(time.RFC3339),
		RequestID: c.GetString("request_id"),
	}
	c.JSON(http.StatusOK, response)
}

// Error 错误响应
func Error(c *gin.Context, code int, message string) {
	response := Response{
		Code:      code,
		Message:   message,
		Timestamp: time.Now().Format(time.RFC3339),
		RequestID: c.GetString("request_id"),
	}
	
	httpStatus := http.StatusOK
	switch code {
	case INVALID_PARAMS:
		httpStatus = http.StatusBadRequest
	case UNAUTHORIZED:
		httpStatus = http.StatusUnauthorized
	case FORBIDDEN:
		httpStatus = http.StatusForbidden
	case NOT_FOUND:
		httpStatus = http.StatusNotFound
	case METHOD_NOT_ALLOWED:
		httpStatus = http.StatusMethodNotAllowed
	case TOO_MANY_REQUESTS:
		httpStatus = http.StatusTooManyRequests
	case INTERNAL_SERVER_ERROR:
		httpStatus = http.StatusInternalServerError
	case SERVICE_UNAVAILABLE:
		httpStatus = http.StatusServiceUnavailable
	}
	
	c.JSON(httpStatus, response)
}

// InvalidParams 参数错误响应
func InvalidParams(c *gin.Context, message string) {
	if message == "" {
		message = MSG_INVALID_PARAMS
	}
	Error(c, INVALID_PARAMS, message)
}

// Unauthorized 未授权响应
func Unauthorized(c *gin.Context, message string) {
	if message == "" {
		message = MSG_UNAUTHORIZED
	}
	Error(c, UNAUTHORIZED, message)
}

// Forbidden 禁止访问响应
func Forbidden(c *gin.Context, message string) {
	if message == "" {
		message = MSG_FORBIDDEN
	}
	Error(c, FORBIDDEN, message)
}

// NotFound 未找到响应
func NotFound(c *gin.Context, message string) {
	if message == "" {
		message = MSG_NOT_FOUND
	}
	Error(c, NOT_FOUND, message)
}

// InternalServerError 服务器内部错误响应
func InternalServerError(c *gin.Context, message string) {
	if message == "" {
		message = MSG_INTERNAL_SERVER_ERROR
	}
	Error(c, INTERNAL_SERVER_ERROR, message)
}

// PageData 分页数据结构体
type PageData struct {
	Items      interface{} `json:"items"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

// NewPageData 创建分页数据
func NewPageData(items interface{}, total int64, page, pageSize int) *PageData {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	
	return &PageData{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}

// SuccessWithPage 分页成功响应
func SuccessWithPage(c *gin.Context, items interface{}, total int64, page, pageSize int) {
	pageData := NewPageData(items, total, page, pageSize)
	Success(c, pageData)
}