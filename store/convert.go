package store

import (
	"time"

	"forum/models"
)

// 本层是做"翻译"：数据库层(gormUser 等) → 对外层(models.User 等)
// 差异：时间字段库里是 time.Time、对外是字符串；数字库里是 int64、对外是 int

func toUser(g *gormUser) *models.User {
	return &models.User{
		ID:       g.ID,
		Username: g.Username,
		Name:     g.Name,
		Role:     g.Role,
		Password: g.Password, // 保留哈希，登录比对要用
	}
}

func toPost(g *gormPost) *models.Post {
	return &models.Post{
		ID:           g.ID,
		Content:      g.Content,
		Author:       toUser(&g.Author),
		LikeCount:    int(g.LikeCount),
		CommentCount: int(g.CommentCount),
		CreatedAt:    g.CreatedAt.Format(time.RFC3339),
	}
}

func toComment(g *gormComment) *models.Comment {
	return &models.Comment{
		ID:        g.ID,
		PostID:    g.PostID,
		Content:   g.Content,
		Author:    toUser(&g.Author),
		CreatedAt: g.CreatedAt.Format(time.RFC3339),
	}
}
