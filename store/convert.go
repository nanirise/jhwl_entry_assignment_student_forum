package store

import (
	"time"

	"forum/models"
)

// ---- 这一层是"架桥"：把数据库层(gormPost等)转成对外层(models.Post等) ----
//
// 数据库卡的字段和对外卡不完全一样：
//   时间字段——库里是 time.Time 对象，对外是 RFC3339 字符串
//   数字类型——库里可能是 int64，对外是 int
// 所以每次从库查出"数据库卡"后，都要经这里"翻译"成 handler 认得的"商品卡"。

// toUser 把 gormUser(数据库层) 转成 models.User(对外层)
func toUser(g *gormUser) *models.User {
	return &models.User{
		ID:       g.ID,
		Username: g.Username,
		Name:     g.Name,
		Role:     g.Role,
		Password: g.Password, // 密码哈希保留，登录时比对要用
	}
}

// toPost 把 gormPost 转成 models.Post
func toPost(g *gormPost) *models.Post {
	return &models.Post{
		ID:           g.ID,
		Content:      g.Content,
		Author:       toUser(&g.Author), // 作者是关联的，一起转出来
		LikeCount:    int(g.LikeCount),
		CommentCount: int(g.CommentCount),
		CreatedAt:    g.CreatedAt.Format(time.RFC3339), // time.Time → 字符串
	}
}

// toComment 把 gormComment 转成 models.Comment
func toComment(g *gormComment) *models.Comment {
	return &models.Comment{
		ID:        g.ID,
		PostID:    g.PostID,
		Content:   g.Content,
		Author:    toUser(&g.Author),
		CreatedAt: g.CreatedAt.Format(time.RFC3339),
	}
}
