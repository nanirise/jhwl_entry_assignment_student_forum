package store

import (
	"time"

	"gorm.io/gorm"
)

// gormUser 用户表
type gormUser struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	Username  string    `gorm:"size:32;uniqueIndex"` // 学号唯一
	Name      string    `gorm:"size:32"`
	Password  string    `gorm:"size:255"` // 存的是哈希
	Role      string    `gorm:"size:10;default:student"`
	CreatedAt time.Time // Gorm 自动填写
	UpdatedAt time.Time
}

func (gormUser) TableName() string { return "users" }

// gormPost 帖子表：AuthorID 存作者ID，Author 为关联用户，查询时可一次带出
type gormPost struct {
	ID           int64  `gorm:"primaryKey;autoIncrement"`
	Content      string `gorm:"type:text"`
	LikeCount    int64
	CommentCount int64
	AuthorID     int64
	Author       gormUser `gorm:"foreignKey:AuthorID"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt // 软删除：只打标记不真删
}

func (gormPost) TableName() string { return "posts" }

// gormComment 评论表
type gormComment struct {
	ID        int64 `gorm:"primaryKey;autoIncrement"`
	PostID    int64
	Content   string `gorm:"type:text"`
	AuthorID  int64
	Author    gormUser `gorm:"foreignKey:AuthorID"`
	CreatedAt time.Time
}

func (gormComment) TableName() string { return "comments" }

// gormLike 点赞表：(user_id, post_id) 复合唯一，数据库层面防重复点赞
type gormLike struct {
	ID        int64 `gorm:"primaryKey;autoIncrement"`
	UserID    int64 `gorm:"uniqueIndex:idx_user_post"`
	PostID    int64 `gorm:"uniqueIndex:idx_user_post"`
	CreatedAt time.Time
}

func (gormLike) TableName() string { return "likes" }
