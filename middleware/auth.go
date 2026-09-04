package middleware

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"

	"forum/pkg/jwt"
	"forum/response"
)

// Auth 鉴权中间件：挡在受保护接口前面的"安检门"
// 每个请求来这里，先验 JWT 手环。验过 → 放行并写入用户信息；验不过 → 返回 401
func Auth() gin.HandlerFunc {
	// 返回一个符合 Gin 中间件签名的函数（记账号、验手环、放行/拦截）
	return func(c *gin.Context) {
		// 1. 从请求头读取 Authorization 字段
		// 格式：Authorization: Bearer <token>
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, response.StatusUnauthorized, "未登录或令牌无效")
			c.Abort() // 中断请求，后面的接口函数不再执行
			return
		}

		// 2. 取出 Bearer 后面的 token 字符串
		// strings.Split 把 "Bearer xxx" 按空格切成 ["Bearer", "xxx"]，取第二段
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, response.StatusUnauthorized, "未登录或令牌无效")
			c.Abort()
			return
		}
		tokenString := parts[1]

		// 3. 用 ParseToken 验证手环真伪与是否过期
		userID, role, err := jwt.ParseToken(tokenString)
		if err != nil {
			// 记录拒绝原因（过期/伪造/格式不对），方便排查"为什么是 401"
			log.Println("鉴权-令牌校验失败:", err)
			response.Error(c, response.StatusUnauthorized, "未登录或令牌无效")
			c.Abort()
			return
		}

		// 4. 把用户信息存进请求上下文，供后续接口使用
		c.Set("userID", userID) // 用户 ID
		c.Set("role", role)     // 角色

		// 5. 全部通过 → 放行，让请求继续往后走到真正的接口函数
		c.Next()
	}
}
