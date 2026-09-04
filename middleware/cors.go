package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORS 解决跨域：网页版 Apifox 在浏览器里跑，请求会触发浏览器的跨域限制。
// 这里统一在响应头里告诉浏览器"我允许你访问"，并处理它先发来的 OPTIONS 预检请求。
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 允许任意来源访问。咱们用 Authorization(Bearer token) 来鉴权，不依赖 Cookie，
		// 所以用 *（全部来源）是安全的；若以后要支持带 Cookie 的跨域，就不能用 *。
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		// 允许前端请求时带的这两个头：请求体格式 + 登录令牌
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// 浏览器在真正发请求前，会先发一个 OPTIONS 预检试探一下。
		// 这里直接回 204(无内容)并中断，不再进业务逻辑；业务接口本身不处理 OPTIONS。
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
