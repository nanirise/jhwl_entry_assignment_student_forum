package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"forum/models"
	"forum/response"
	"forum/store"
)

// LikeHandler 处理点赞相关接口：⑧ 点赞/取消，⑩ 批量查状态
type LikeHandler struct {
	store *store.Store
}

// NewLikeHandler 创建点赞处理器，注入仓库
func NewLikeHandler(s *store.Store) *LikeHandler {
	return &LikeHandler{store: s}
}

// currentUser 从 token 取当前登录用户
func (h *LikeHandler) currentUser(c *gin.Context) *models.User {
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

// Like 点赞/取消：POST /api/v1/posts/:post_id/like，点一次=赞，再点=取消
func (h *LikeHandler) Like(c *gin.Context) {
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

	isLiked, err := h.store.ToggleLike(user.ID, postID)
	if err != nil {
		response.Error(c, response.StatusNotFound, "帖子不存在")
		return
	}

	response.Success(c, response.StatusOK, models.LikeToggleResult{
		PostID:  postID,
		IsLiked: isLiked,
	})
}

// GetLikeStatus 批量点赞状态：POST /api/v1/posts/likes，请求体 {"post_ids":[...]}
func (h *LikeHandler) GetLikeStatus(c *gin.Context) {
	user := h.currentUser(c)
	if user == nil {
		response.Error(c, response.StatusUnauthorized, "未登录或令牌无效")
		return
	}

	var req models.LikeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.StatusBadRequest, "参数校验失败")
		return
	}

	// 校验：不能为空、上限 100、id 必须为正整数
	if len(req.PostIDs) == 0 || len(req.PostIDs) > 100 {
		response.Error(c, response.StatusBadRequest, "参数校验失败")
		return
	}
	for _, pid := range req.PostIDs {
		if pid < 1 {
			response.Error(c, response.StatusBadRequest, "参数校验失败")
			return
		}
	}

	statuses := h.store.GetLikeStatus(user.ID, req.PostIDs)
	response.Success(c, response.StatusOK, gin.H{
		"status": statuses,
	})
}
