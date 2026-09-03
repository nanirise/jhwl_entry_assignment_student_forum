package models

// Post 一张帖子的"档案卡"：定义一帖子长什么样
type Post struct {
	ID           int64  `json:"id"`           // 帖子 ID
	Content      string `json:"content"`      // 正文
	Author       *User  `json:"author"`       // 作者（贴上去的小照片，是一个 User）
	LikeCount    int    `json:"like_count"`   // 点赞数（Day6 才真的会数）
	CommentCount int    `json:"comment_count"` // 评论数
	CreatedAt    string `json:"created_at"`   // 发布时间
}

// PostDetail 帖子详情：帖子本身 + 它下面的所有评论
// 详情页既要看帖，又要看评论，所以把它们装进同一个箱子
type PostDetail struct {
	ID           int64     `json:"id"`
	Content      string    `json:"content"`
	Author       *User     `json:"author"`
	LikeCount    int       `json:"like_count"`
	CommentCount int       `json:"comment_count"`
	CreatedAt    string    `json:"created_at"`
	Comments     []*Comment `json:"comments"`
}
