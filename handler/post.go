package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"forum/models"
	"forum/response"
	"forum/store"
)

// PostHandler 处理"帖子与评论"相关的接口
type PostHandler struct {
	store *store.Store
}

// NewPostHandler 创建一个帖子处理器，注入仓库
func NewPostHandler(s *store.Store) *PostHandler {
	return &PostHandler{store: s}
}

// currentUser 从 token 里取出当前登录用户，多个接口复用
func (h *PostHandler) currentUser(c *gin.Context) *models.User {
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

// Create 发布帖子：POST /api/v1/posts
// 作者必须从 JWT 拿，禁止请求体传 user_id（文档要求）
func (h *PostHandler) Create(c *gin.Context) {
	user := h.currentUser(c)
	if user == nil {
		response.Error(c, response.StatusUnauthorized, "未登录或令牌无效")
		return
	}

	var req models.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.StatusBadRequest, "参数校验失败")
		return
	}

	created := h.store.CreatePost(req.Content, user)
	response.Success(c, response.StatusCreated, created)
}

// List 分页获取帖子列表：GET /api/v1/posts
func (h *PostHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if pageSize < 1 {
		pageSize = 20
	} else if pageSize > 100 {
		pageSize = 100
	}

	items, total := h.store.ListPosts(page, pageSize)
	response.Success(c, response.StatusOK, gin.H{
		"items": items,
		"meta": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
	})
}

// Get 获取帖子详情：GET /api/v1/posts/:post_id（带上它的评论）
func (h *PostHandler) Get(c *gin.Context) {
	postID, err := strconv.ParseInt(c.Param("post_id"), 10, 64)
	if err != nil || postID < 1 {
		response.Error(c, response.StatusBadRequest, "参数校验失败")
		return
	}

	post := h.store.GetPostByID(postID)
	if post == nil {
		response.Error(c, response.StatusNotFound, "帖子不存在")
		return
	}

	comments := h.store.GetCommentsByPostID(postID)
	response.Success(c, response.StatusOK, models.PostDetail{
		ID:           post.ID,
		Content:      post.Content,
		Author:       post.Author,
		LikeCount:    post.LikeCount,
		CommentCount: len(comments),
		CreatedAt:    post.CreatedAt,
		Comments:     comments,
	})
}

// CreateComment 发表评论：POST /api/v1/posts/:post_id/comment
func (h *PostHandler) CreateComment(c *gin.Context) {
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

	var req models.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.StatusBadRequest, "参数校验失败")
		return
	}

	created, err := h.store.CreateComment(&models.Comment{
		PostID:  postID,
		Content: req.Content,
		Author:  user,
	})
	if err != nil {
		response.Error(c, response.StatusNotFound, "帖子不存在")
		return
	}
	response.Success(c, response.StatusCreated, created)
}

// Delete 删除自己的帖子：DELETE /api/v1/posts/:post_id
// 只能删"你这个作者"发的帖子；别人的帖子 → 403（删别人的统一走 ⑦ 管理员接口）
func (h *PostHandler) Delete(c *gin.Context) {
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

	// 3. 帖子得存在才谈得上删
	post := h.store.GetPostByID(postID)
	if post == nil {
		response.Error(c, response.StatusNotFound, "帖子不存在")
		return
	}

	// 4. 权限：这个帖子的作者必须是"你"，否则就是删别人的 → 403
	if post.Author.ID != user.ID {
		response.Error(c, response.StatusForbidden, "无权删除他人的帖子")
		return
	}

	// 5. 上面已确认帖子存在，级联删帖不会出错；文档要求成功时 data 为 null
	_ = h.store.DeletePost(postID)
	response.Success(c, response.StatusOK, nil)
}
