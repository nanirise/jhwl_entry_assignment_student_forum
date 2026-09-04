package models

// Post 帖子
type Post struct {
	ID           int64  `json:"id"`
	Content      string `json:"content"`
	Author       *User  `json:"author"`
	LikeCount    int    `json:"like_count"`
	CommentCount int    `json:"comment_count"`
	CreatedAt    string `json:"created_at"`
}

// PostDetail 帖子详情：帖子本身 + 它下面的所有评论
type PostDetail struct {
	ID           int64      `json:"id"`
	Content      string     `json:"content"`
	Author       *User      `json:"author"`
	LikeCount    int        `json:"like_count"`
	CommentCount int        `json:"comment_count"`
	CreatedAt    string     `json:"created_at"`
	Comments     []*Comment `json:"comments"`
}
