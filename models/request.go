package models

// RegisterRequest 注册请求体：顾客递进来的"点菜单"
// 分开定义接收结构体：既能接收密码，JSON 响应时密码又不会泄漏出去
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`                     // 学号/工号，纯数字
	Name     string `json:"name"     binding:"required"`                     // 姓名
	Password string `json:"password" binding:"required"`                     // 密码，8~16位
	Role     string `json:"role"     binding:"required,oneof=student admin"` // 角色：student / admin
}

// LoginRequest 登录请求体：顾客报"账号 + 密码"
// 登录时不需要 role，角色由服务器从库里查出
type LoginRequest struct {
	Username string `json:"username" binding:"required"` // 学号/工号
	Password string `json:"password" binding:"required"` // 密码
}
