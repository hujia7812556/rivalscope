// Package dto 定义接口的请求/响应结构与统一响应封装。
package dto

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构:{code, message, data}。
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Error 业务错误,携带业务码、HTTP 状态码与提示信息。
type Error struct {
	Code       int    // 响应体中的业务码
	HTTPStatus int    // HTTP 状态码
	Message    string // 提示信息
}

// Error 实现 error 接口。
func (e *Error) Error() string { return e.Message }

// NewError 创建业务错误。
func NewError(code, httpStatus int, message string) *Error {
	return &Error{Code: code, HTTPStatus: httpStatus, Message: message}
}

// 预定义错误码:4xx/5xx 与 HTTP 语义一致,业务错误从 1001 起。
var (
	ErrBadRequest          = NewError(400, http.StatusBadRequest, "请求参数错误")
	ErrUnauthorized        = NewError(401, http.StatusUnauthorized, "未登录或登录已过期")
	ErrNotFound            = NewError(404, http.StatusNotFound, "资源不存在")
	ErrInternalServerError = NewError(500, http.StatusInternalServerError, "服务器内部错误")
	ErrLoginFailed         = NewError(1002, http.StatusUnauthorized, "用户名或密码错误")
	ErrWeiboUidExists      = NewError(2001, http.StatusBadRequest, "该 uid 已存在")
)

// Success 输出成功响应;data 为 nil 时输出空对象,保证 data 字段稳定为 {}。
func Success(c *gin.Context, data interface{}) {
	if data == nil {
		data = gin.H{}
	}
	c.JSON(http.StatusOK, Response{Code: 0, Message: "ok", Data: data})
}

// Fail 输出失败响应;err 为 *Error 时使用其业务码与 HTTP 状态码,否则按 500 处理。
func Fail(c *gin.Context, err error) {
	if e, ok := err.(*Error); ok {
		c.JSON(e.HTTPStatus, Response{Code: e.Code, Message: e.Message, Data: gin.H{}})
		return
	}
	c.JSON(http.StatusInternalServerError,
		Response{Code: 500, Message: "服务器内部错误", Data: gin.H{}})
}
