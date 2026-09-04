package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"forum/models"
	"forum/response"
	"forum/store"
)

// AdminHandler 处理管理员专属接口：⑦ 删除任意帖子
type AdminHandler struct {
	store *store.Store
}

// NewAdminHandler 创建管理员处理器，注入仓库
func NewAdminHandler(s *store.Store) *AdminHandler {
	return &AdminHandler{store: s}
}

// currentUser 从 token 取当前登录用户
func (h *AdminHandler) currentUser(c *gin.Context) *models.User {
	userID, ok := c.Get("userID")
	if !ok {
		return nil
	}
	id, ok := userID.(int64)
	if !ok {
		return nil
	}
	return h.store.GetUserByID(id)
}

// DeletePost 管理员删除任意帖子：DELETE /api/v1/admin/posts/:post_id
func (h *AdminHandler) DeletePost(c *gin.Context) {
	user := h.currentUser(c)
	if user == nil {
		response.Error(c, response.StatusUnauthorized, "未登录或令牌无效")
		return
	}

	postID, err := strconv.ParseInt(c.Param("post_id"), 10, 64)
	if err != nil || postID < 1 {
		response.Error(c, response.StatusBadRequest, "参数校验失败")
		return
	}

	// 身份校验：必须是 admin，否则 403
	if user.Role != "admin" {
		response.Error(c, response.StatusForbidden, "仅管理员可删除任意帖子")
		return
	}

	if err := h.store.DeletePost(postID); err != nil {
		response.Error(c, response.StatusNotFound, "帖子不存在")
		return
	}

	response.Success(c, response.StatusOK, nil)
}
