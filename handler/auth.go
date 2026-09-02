package handler

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"forum/models"
	"forum/pkg/jwt"
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

// Login 处理 POST /api/v1/auth/login 登录请求
func (h *AuthHandler) Login(c *gin.Context) {
	// 1. 解析请求体：顾客报"账号 + 密码"
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.StatusBadRequest, "参数校验失败")
		return
	}

	// 2. 按学号查人：查不到说明账号不存在
	user := h.store.GetUserByUsername(req.Username)
	if user == nil {
		// 统一报"账号或密码错误"，不区分哪个错（防信息泄露），返回 401
		response.Error(c, response.StatusUnauthorized, "账号或密码错误")
		return
	}

	// 3. 验证密码：CompareHashAndPassword 会拿输入的明文重新哈希，跟库里的哈希比对
	// 如果对不上，返回 err（账号或密码错误）
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		response.Error(c, response.StatusUnauthorized, "账号或密码错误")
		return
	}

	// 4. 密码正确，签发 JWT 手环：用用户 ID、学号、角色去生成
	tokenString, expiresIn, err := jwt.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		response.Error(c, response.StatusInternalServerError, "令牌生成失败")
		return
	}

	// 5. 组装登录成功响应（文档要求的四样：access_token/token_type/expires_in/user）
	loginResp := models.LoginResponse{
		AccessToken: tokenString,
		TokenType:   "Bearer",
		ExpiresIn:   expiresIn,
		User:        *user,
	}
	response.Success(c, response.StatusOK, loginResp)
}

// Me 处理 GET /api/v1/auth/me：返回当前登录用户的信息
// 这是一个"受保护接口"的示例：必须经过鉴权中间件才能访问
func (h *AuthHandler) Me(c *gin.Context) {
	// 中间件已经把用户 ID 存进了上下文，这里取出来
	userID, _ := c.Get("userID")

	// 用 ID 查用户
	user := h.store.GetUserByID(userID.(int64))
	if user == nil {
		response.Error(c, response.StatusNotFound, "用户不存在")
		return
	}

	response.Success(c, response.StatusOK, user)
}
