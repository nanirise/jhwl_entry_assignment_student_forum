package models

// RegisterRequest 注册请求体：接收并校验注册字段
type RegisterRequest struct {
	Username string `json:"username" binding:"required,numeric,max=32"`      // 学号，纯数字，≤32 位
	Name     string `json:"name"     binding:"required,max=32"`              // 姓名，≤32 位
	Password string `json:"password" binding:"required,min=8,max=16"`        // 密码，8~16 位
	Role     string `json:"role"     binding:"required,oneof=student admin"` // 角色
}

// LoginRequest 登录请求体
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// CreatePostRequest 发帖请求体：只收 content，作者由服务端从 token 取
type CreatePostRequest struct {
	Content string `json:"content" binding:"required,min=1,max=2000"`
}

// CreateCommentRequest 评论请求体
type CreateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1,max=1000"`
}
