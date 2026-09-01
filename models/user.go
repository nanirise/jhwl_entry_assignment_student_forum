package models

// User 用户档案卡：定义"一个用户长什么样"
type User struct {
	ID       int64  `json:"id"`       // 用户 ID
	Username string `json:"username"` // 学号或工号，纯数字
	Name     string `json:"name"`     // 姓名
	Role     string `json:"role"`     // 角色：student / admin
	Password string `json:"-"`        // 密码（哈希后），json:"-" 表示永远不输出到 JSON
}
