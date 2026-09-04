package models

// User 用户档案
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"` // 学号/工号，纯数字
	Name     string `json:"name"`     // 姓名
	Role     string `json:"role"`     // student / admin
	Password string `json:"-"`        // 密码(哈希)：json:"-" 保证永不输出到响应
}
