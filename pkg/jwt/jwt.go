package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 密钥：给 JWT 签名的钢印，只有服务器知道。真实项目应放配置文件/环境变量，不可写死/泄露
var secret = []byte("forum-secret-key-2026") // [假设] 正式版用 Viper 配置，暂用固定值

// GenerateToken 签发 JWT：入参用户ID、学号、角色，返回 token 字符串和有效期(秒)
func GenerateToken(userID int64, username, role string) (string, int64, error) {
	expireSeconds := int64(7200) // 2 小时

	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(time.Second * time.Duration(expireSeconds)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", 0, err
	}

	return tokenString, expireSeconds, nil
}

// ParseToken 校验 JWT：返回用户ID、角色；无效/过期/伪造则返回 error
func ParseToken(tokenString string) (int64, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		// 只认 HS256，防止攻击者换算法绕过
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return secret, nil
	})
	if err != nil {
		return 0, "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, "", jwt.ErrTokenInvalidClaims
	}

	// JWT 里的数字默认解析成 float64，需转回整数
	userID, _ := claims["user_id"].(float64)
	role, _ := claims["role"].(string)

	return int64(userID), role, nil
}
