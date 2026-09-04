package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORS 解决跨域：网页版 Apifox 在浏览器里跑，请求会触发浏览器的"同源限制"。
// 关键在"预检谈判"：浏览器真正请求前会先发一个 OPTIONS，列出它这次想带的请求头，
// 服务器必须回"你带的这些我全都放行"；否则浏览器就判 CORS 拒绝。
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 允许的来源：把浏览器带来的 Origin 原样回给它（而不是写死 *）。
		//    无论 Apifox 网页从哪个地址来，都能被放行；同时放行带凭据(如 Cookie)的请求。
		origin := c.GetHeader("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Vary", "Origin")
		} else {
			// 没带 Origin（比如 curl 直连）时，就用 * 兜底。
			c.Header("Access-Control-Allow-Origin", "*")
		}

		// 2. 允许的请求方法
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// 3. 允许的请求头：核心修复！不再猜它用哪几个头——预检里浏览器申请了哪些，
		//    就原样回哪些，"你带什么我就放行什么"，绝不会因为漏了某个自定义头而拦截。
		reqHeaders := c.GetHeader("Access-Control-Request-Headers")
		if reqHeaders != "" {
			c.Header("Access-Control-Allow-Headers", reqHeaders)
		} else {
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		}

		// 预检请求 OPTIONS 本身：回 204(无内容)并中断，不再进业务逻辑。
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
