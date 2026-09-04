package models

// Comment 评论
type Comment struct {
	ID        int64  `json:"id"`
	PostID    int64  `json:"post_id"`
	Content   string `json:"content"`
	Author    *User  `json:"author"`
	CreatedAt string `json:"created_at"`
}
