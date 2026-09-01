package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Body 统一响应结构：所有接口的返回都装进这个"标准纸箱"
type Body struct {
	Code int    `json:"code"` // 业务状态码：0=成功，出错时=HTTP状态码
	Msg  string `json:"msg"`  // 面向调用方的提示
	Data any    `json:"data"` // 真正的数据（任意类型）
}

// Success 成功响应：code=0，msg="success"，data 为返回的数据
// 注意：业务成功不代表 HTTP 一定是 200，有些接口文档要求 201（如注册、发帖、评论）
func Success(c *gin.Context, httpStatus int, data any) {
	c.JSON(httpStatus, Body{
		Code: 0,
		Msg:  "success",
		Data: data,
	})
}

// Error 失败响应：code=HTTP状态码，msg 为错误描述，data 固定为 nil
func Error(c *gin.Context, httpStatus int, msg string) {
	c.JSON(httpStatus, Body{
		Code: httpStatus,
		Msg:  msg,
		Data: nil,
	})
}

// 顺手替你把常用的状态码常量准备好，避免在代码里写魔法数字
// 这些常量来自 net/http，语义清晰、面试也显得专业
var (
	StatusOK    = http.StatusOK    // 200 成功
	StatusCreated = http.StatusCreated // 201 创建成功
	StatusBadRequest   = http.StatusBadRequest   // 400 参数有误
	StatusUnauthorized = http.StatusUnauthorized // 401 未登录/令牌无效
	StatusForbidden    = http.StatusForbidden    // 403 无权操作
	StatusNotFound     = http.StatusNotFound     // 404 资源不存在
	StatusConflict     = http.StatusConflict     // 409 冲突（如用户名已存在）
	StatusInternalServerError = http.StatusInternalServerError // 500 服务器内部错误
)
