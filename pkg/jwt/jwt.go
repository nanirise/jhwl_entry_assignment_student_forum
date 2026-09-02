package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 密钥（secret）：给 JWT 签名的"钢印模板"
// 只有服务器知道，顾客改不了、伪造不了。真实项目里应放在配置文件/环境变量，绝不写死/泄露！
var secret = []byte("forum-secret-key-2026") // [假设] 正式版用 Viper 配置，暂用固定值

// GenerateToken 根据用户信息签发一个 JWT 手环
// 返回签名后的 token 字符串、以及有效期（秒）
// 入参：用户ID、学号、角色 —— 信息记录在手环里，服务器之后靠它识别是谁
func GenerateToken(userID int64, username, role string) (string, int64, error) {
	// 设定有效期 7200 秒 = 2 小时（文档要求 expires_in=7200）
	expireSeconds := int64(7200)

	// 构造 JWT 的"负载(claims)"：手环上记录的信息 + 过期时间
	claims := jwt.MapClaims{
		"user_id":  userID,   // 用户 ID
		"username": username, // 学号
		"role":     role,     // 角色
		"exp":      time.Now().Add(time.Second * time.Duration(expireSeconds)).Unix(), // 过期时间（Unix毫秒）
	}

	// 用 HS256 签名算法 + 密钥，把负载"盖钢印"生成手环
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", 0, err
	}

	return tokenString, expireSeconds, nil
}

// ParseToken 验证一个 JWT 手环
// 返回手环上记录的用户 ID 和角色；若无效/过期/伪造，返回 error
// 保安验手环：能验过且未过期 → 放行；否则拒绝
func ParseToken(tokenString string) (int64, string, error) {
	// 解析 token，提供密钥来验证签名
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		// 检查签名算法是不是我们期望的 HS256（防止攻击者换算法）
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return secret, nil
	})
	if err != nil {
		return 0, "", err
	}

	// 取出负载里的信息（claims）
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, "", jwt.ErrTokenInvalidClaims
	}

	// 读取 user_id 和 role（断言类型为 float64，因为 JSON 数字默认是 float64）
	userID, _ := claims["user_id"].(float64)
	role, _ := claims["role"].(string)

	return int64(userID), role, nil
}
