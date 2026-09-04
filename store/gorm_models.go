package store

import (
	"time"

	"gorm.io/gorm"
)

// gormUser 用户表：这张"入库登记卡"告诉 Gorm 怎么建 users 表
type gormUser struct {
	ID        int64  `gorm:"primaryKey;autoIncrement"` // 主键，自增
	Username  string `gorm:"size:32;uniqueIndex"`      // 唯一：学号不能重复
	Name      string `gorm:"size:32"`
	Password  string `gorm:"size:255"` // 存的是哈希后的密码
	Role      string `gorm:"size:10;default:student"`
	CreatedAt time.Time // Gorm 自动在该列写"创建时间"
	UpdatedAt time.Time // Gorm 自动写"更新时间"
}

// TableName 显式指定表名（Gorm 默认会复数成 users，这里保险写清楚）
func (gormUser) TableName() string { return "users" }

// gormPost 帖子表：作者用 AuthorID 记，Author 是"关联的用户"，查询时可一次带出
type gormPost struct {
	ID           int64  `gorm:"primaryKey;autoIncrement"`
	Content      string `gorm:"type:text"` // 长文
	LikeCount    int64
	CommentCount int64
	AuthorID     int64
	Author       gormUser     `gorm:"foreignKey:AuthorID"` // 作者关联
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt // 软删除：删帖只是打标记，不真删
}

func (gormPost) TableName() string { return "posts" }

// gormComment 评论表：也是"挂在某个人身上"，AuthorID + Author 关联
type gormComment struct {
	ID        int64  `gorm:"primaryKey;autoIncrement"`
	PostID    int64
	Content   string `gorm:"type:text"`
	AuthorID  int64
	Author    gormUser `gorm:"foreignKey:AuthorID"`
	CreatedAt time.Time
}

func (gormComment) TableName() string { return "comments" }

// gormLike 点赞表：同一用户同一帖子只能赞一次（用复合唯一锁住，防重复点赞）
type gormLike struct {
	ID        int64 `gorm:"primaryKey;autoIncrement"`
	UserID    int64 `gorm:"uniqueIndex:idx_user_post"` // 两个字段一起组成唯一索引
	PostID    int64 `gorm:"uniqueIndex:idx_user_post"`
	CreatedAt time.Time
}

func (gormLike) TableName() string { return "likes" }
