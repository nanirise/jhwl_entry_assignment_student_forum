package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"forum/models"
	"forum/response"
	"forum/store"
)

// AdminHandler 处理"管理员专属"的接口：⑦ 删除任意帖子
type AdminHandler struct {
	store *store.Store
}

// NewAdminHandler 创建一个管理员处理器，注入仓库
func NewAdminHandler(s *store.Store) *AdminHandler {
	return &AdminHandler{store: s}
}

// currentUser 从 token 里取出当前登录用户（用法和前面一致）
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
// 不管帖子是谁发的，只要你是 admin 就能删
func (h *AdminHandler) DeletePost(c *gin.Context) {
	// 1. 没登录不让删
	user := h.currentUser(c)
	if user == nil {
		response.Error(c, response.StatusUnauthorized, "未登录或令牌无效")
		return
	}

	// 2. 路径里的帖子 ID 必须是合法数字
	postID, err := strconv.ParseInt(c.Param("post_id"), 10, 64)
	if err != nil || postID < 1 {
		response.Error(c, response.StatusBadRequest, "参数校验失败")
		return
	}

	// 3. 身份校验：你的角色必须是 admin，否则 → 403
	if user.Role != "admin" {
		response.Error(c, response.StatusForbidden, "仅管理员可删除任意帖子")
		return
	}

	// 4. 帖子得存在，删除才不落空
	if err := h.store.DeletePost(postID); err != nil {
		response.Error(c, response.StatusNotFound, "帖子不存在")
		return
	}

	// 5. 成功，文档要求 data 为 null
	response.Success(c, response.StatusOK, nil)
}
