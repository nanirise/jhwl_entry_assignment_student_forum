package models

// Comment 一条评论的"档案卡"
type Comment struct {
	ID        int64  `json:"id"`         // 评论 ID
	PostID    int64  `json:"post_id"`    // 挂在哪个帖子下
	Content   string `json:"content"`    // 评论内容
	Author    *User  `json:"author"`     // 评论者
	CreatedAt string `json:"created_at"` // 评论时间
}
