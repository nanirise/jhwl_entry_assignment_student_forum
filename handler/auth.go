package handler

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"forum/models"
	"forum/response"
	"forum/store"
)

// AuthHandler 处理"用户与鉴权"相关的接口
type AuthHandler struct {
	store *store.Store // 持有仓库，方便存取数据
}

// NewAuthHandler 创建一个鉴权处理器，把仓库注入进来
func NewAuthHandler(s *store.Store) *AuthHandler {
	return &AuthHandler{store: s}
}

// Register 处理 POST /api/v1/auth/register 注册请求
func (h *AuthHandler) Register(c *gin.Context) {
	// 1. 解析请求体：把顾客递来的"点菜单"填进 RegisterRequest 结构体
	var req models.RegisterRequest
	// ShouldBindJSON 会自动校验 binding:"required"，没填会返回错误
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.StatusBadRequest, "参数校验失败")
		return
	}

	// 2. 密码哈希：把明文密码变成不可逆的乱码（这是文档硬性要求）
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.Error(c, response.StatusInternalServerError, "密码处理失败")
		return
	}

	// 3. 组装 User（注意：存的是哈希后的密码，不是明文）
	user := &models.User{
		Username: req.Username,
		Name:     req.Name,
		Password: string(hashed), // 关键：这里放哈希值，绝不放明文
		Role:     req.Role,
	}

	// 4. 存入存储，若用户名已存在则返回冲突(409)
	created, err := h.store.CreateUser(user)
	if err != nil {
		response.Error(c, response.StatusConflict, "用户名已存在")
		return
	}

	// 5. 注册成功，返回 201 + 用户信息（User 结构体不含密码字段，所以不会泄漏）
	response.Success(c, response.StatusCreated, created)
}
