package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Body 统一响应结构：所有接口都装进这个标准格式
type Body struct {
	Code int    `json:"code"` // 0=成功，出错时=HTTP 状态码
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

// Success 成功响应：code=0，msg=success。注意 HTTP 状态码可能因文档要求是 201
func Success(c *gin.Context, httpStatus int, data any) {
	c.JSON(httpStatus, Body{
		Code: 0,
		Msg:  "success",
		Data: data,
	})
}

// Error 失败响应：code=HTTP 状态码，data 固定 nil
func Error(c *gin.Context, httpStatus int, msg string) {
	c.JSON(httpStatus, Body{
		Code: httpStatus,
		Msg:  msg,
		Data: nil,
	})
}

// 常用状态码常量，避免在代码里写死魔法数字
var (
	StatusOK                  = http.StatusOK                  // 200
	StatusCreated             = http.StatusCreated             // 201
	StatusBadRequest          = http.StatusBadRequest          // 400
	StatusUnauthorized        = http.StatusUnauthorized        // 401
	StatusForbidden           = http.StatusForbidden           // 403
	StatusNotFound            = http.StatusNotFound            // 404
	StatusConflict            = http.StatusConflict            // 409
	StatusInternalServerError = http.StatusInternalServerError // 500
)
