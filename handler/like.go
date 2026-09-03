package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"forum/models"
	"forum/response"
	"forum/store"
)

// LikeHandler 处理"点赞"相关的接口：⑧ 点赞/取消，⑩ 批量查点赞状态
type LikeHandler struct {
	store *store.Store
}

// NewLikeHandler 创建一个点赞处理器，注入仓库
func NewLikeHandler(s *store.Store) *LikeHandler {
	return &LikeHandler{store: s}
}

// currentUser 从 token 里取出当前登录用户（和 PostHandler 里的用法一致）
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

// Like 点赞/取消点赞：POST /api/v1/posts/:post_id/like
// 点一下=赞，再点=取消（开关）。返回切换后的当前状态。
func (h *LikeHandler) Like(c *gin.Context) {
	// 1. 没登录就不让点
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

	// 3. 交给仓库"取反 + 更新赞数"，帖子不存在它会给错误
	isLiked, err := h.store.ToggleLike(user.ID, postID)
	if err != nil {
		response.Error(c, response.StatusNotFound, "帖子不存在")
		return
	}

	// 4. 返回切换后的状态
	response.Success(c, response.StatusOK, models.LikeToggleResult{
		PostID:  postID,
		IsLiked: isLiked,
	})
}

// GetLikeStatus 批量查询当前用户对多个帖子的点赞状态：POST /api/v1/posts/likes
// 请求体：{"post_ids":[20001,20002,...]}，一次性查一堆帖子你赞没赞
func (h *LikeHandler) GetLikeStatus(c *gin.Context) {
	// ① 没登录就不让查
	user := h.currentUser(c)
	if user == nil {
		response.Error(c, response.StatusUnauthorized, "未登录或令牌无效")
		return
	}

	// ② 读请求体；post_ids 缺失或不是数组会让 binding="required" 报错 → 400
	var req models.LikeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.StatusBadRequest, "参数校验失败")
		return
	}

	// ③ 校验：post_ids 为空、超过上限(100)、或单个 ID 不合法，都算参数非法 → 400
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

	// ④ 批量查状态（脏活交给仓库）
	statuses := h.store.GetLikeStatus(user.ID, req.PostIDs)

	// ⑤ 按文档 200 定义，包一层 status
	response.Success(c, response.StatusOK, gin.H{
		"status": statuses,
	})
}
