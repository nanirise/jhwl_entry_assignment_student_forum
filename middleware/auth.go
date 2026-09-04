package middleware

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"

	"forum/pkg/jwt"
	"forum/response"
)

// Auth 鉴权中间件：挡在受保护接口前的"安检门"，验 JWT，通过则放行并写入用户信息
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 取请求头 Authorization，格式：Bearer <token>
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, response.StatusUnauthorized, "未登录或令牌无效")
			c.Abort()
			return
		}

		// 2. 拆出 Bearer 后面的 token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, response.StatusUnauthorized, "未登录或令牌无效")
			c.Abort()
			return
		}
		tokenString := parts[1]

		// 3. 验签 + 查过期
		userID, role, err := jwt.ParseToken(tokenString)
		if err != nil {
			log.Println("鉴权-令牌校验失败:", err) // 记日志，方便排查 401 的具体原因
			response.Error(c, response.StatusUnauthorized, "未登录或令牌无效")
			c.Abort()
			return
		}

		// 4. 把身份放进上下文，供后续接口使用
		c.Set("userID", userID)
		c.Set("role", role)

		// 5. 放行
		c.Next()
	}
}
