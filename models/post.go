package models

type Post struct {
	ID           int64  `json:"id"`
	Content      string `json:"content"`
	LikeCount    int    `json:"like_count"`
	CommentCount int    `json:"comment_count"`
	CreatedAt    string `json:"created_at"`
}
