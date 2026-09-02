package models

// LoginResponse 登录成功后的返回内容
// 文档要求：access_token / token_type / expires_in / user
type LoginResponse struct {
	AccessToken string `json:"access_token"` // JWT 访问令牌（手环）
	TokenType   string `json:"token_type"`   // 固定为 "Bearer"
	ExpiresIn   int64  `json:"expires_in"`   // 令牌有效期（秒），文档是 7200
	User        User   `json:"user"`         // 用户信息（不含密码）
}
