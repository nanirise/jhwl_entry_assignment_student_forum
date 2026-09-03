package models

// LikeToggleResult ⑧ 点赞接口的返回：当前这个帖子的点赞状态
// 注意字段叫 is_liked（跟 ⑩ 的 liked 不同）
type LikeToggleResult struct {
	PostID  int64 `json:"post_id"`
	IsLiked bool  `json:"is_liked"`
}

// LikeStatus ⑩ 批量点赞状态里的一项
type LikeStatus struct {
	PostID int64 `json:"post_id"`
	Liked  bool  `json:"liked"`
}

// LikeRequest ⑩ 批量查询点赞状态的请求体
type LikeRequest struct {
	PostIDs []int64 `json:"post_ids" binding:"required"`
}
