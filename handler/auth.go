package handler

import (
	"log"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"forum/models"
	"forum/pkg/jwt"
	"forum/response"
	"forum/store"
)

// AuthHandler 处理用户与鉴权相关接口
type AuthHandler struct {
	store *store.Store
}

// NewAuthHandler 创建鉴权处理器，注入仓库
func NewAuthHandler(s *store.Store) *AuthHandler {
	return &AuthHandler{store: s}
}

// Register 注册：POST /api/v1/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.StatusBadRequest, "参数校验失败")
		return
	}

	// 密码哈希：明文变不可逆乱码（文档硬性要求）
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("注册-密码哈希失败:", err)
		response.Error(c, response.StatusInternalServerError, "密码处理失败")
		return
	}

	// 存的是哈希，不是明文
	user := &models.User{
		Username: req.Username,
		Name:     req.Name,
		Password: string(hashed),
		Role:     req.Role,
	}

	created, err := h.store.CreateUser(user)
	if err != nil {
		log.Println("注册-创建用户失败:", err) // 记原因，区分重名与数据库故障
		response.Error(c, response.StatusConflict, "用户名已存在")
		return
	}

	response.Success(c, response.StatusCreated, created)
}

// Login 登录：POST /api/v1/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.StatusBadRequest, "参数校验失败")
		return
	}

	user := h.store.GetUserByUsername(req.Username)
	if user == nil {
		// 统一报"账号或密码错误"，不暴露到底哪个错
		response.Error(c, response.StatusUnauthorized, "账号或密码错误")
		return
	}

	// 密码比对：重新哈希输入并与库中哈希比对
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		response.Error(c, response.StatusUnauthorized, "账号或密码错误")
		return
	}

	tokenString, expiresIn, err := jwt.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		log.Println("登录-生成令牌失败:", err)
		response.Error(c, response.StatusInternalServerError, "令牌生成失败")
		return
	}

	loginResp := models.LoginResponse{
		AccessToken: tokenString,
		TokenType:   "Bearer",
		ExpiresIn:   expiresIn,
		User:        *user,
	}
	response.Success(c, response.StatusOK, loginResp)
}

// Me 当前登录用户信息：GET /api/v1/auth/me
func (h *AuthHandler) Me(c *gin.Context) {
	userID, _ := c.Get("userID")
	user := h.store.GetUserByID(userID.(int64))
	if user == nil {
		response.Error(c, response.StatusNotFound, "用户不存在")
		return
	}

	response.Success(c, response.StatusOK, user)
}
