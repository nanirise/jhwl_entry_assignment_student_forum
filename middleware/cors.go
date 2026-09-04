package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORS 解决跨域：网页版 Apifox 在浏览器里跑，会触发同源限制。
// 关键在预检(OPTIONS)：浏览器先列出它想带的请求头，服务器必须放行，否则判 CORS 拒绝
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 允许的来源：把浏览器带来的 Origin 原样回给它（并放行携带凭据的请求）
		origin := c.GetHeader("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Vary", "Origin")
		} else {
			c.Header("Access-Control-Allow-Origin", "*")
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// 允许的请求头：回显浏览器申请的头，不漏掉自定义头（核心修复）
		reqHeaders := c.GetHeader("Access-Control-Request-Headers")
		if reqHeaders != "" {
			c.Header("Access-Control-Allow-Headers", reqHeaders)
		} else {
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		}

		// 预检请求本身：回 204 并中断，不再进业务逻辑
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
