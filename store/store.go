package store

import (
	"fmt"
	"log"

	"forum/models"
	"gorm.io/gorm"
)

// Store 仓储：用 MySQL + Gorm 查存数据，对外方法名/参数/返回值不变
type Store struct {
	db *gorm.DB
}

// New 创建仓储：传入已连库的 *gorm.DB
func New(db *gorm.DB) *Store {
	return &Store{db: db}
}

// GetUserByUsername 按学号查用户；查不到返回 nil
func (s *Store) GetUserByUsername(username string) *models.User {
	var g gormUser
	if err := s.db.Where("username = ?", username).First(&g).Error; err != nil {
		return nil
	}
	return toUser(&g)
}

// GetUserByID 按用户ID查用户；查不到返回 nil
func (s *Store) GetUserByID(id int64) *models.User {
	var g gormUser
	if err := s.db.First(&g, id).Error; err != nil {
		return nil
	}
	return toUser(&g)
}

// CreateUser 创建用户；用户名已存在返回错误（调用方据此返回 409）
func (s *Store) CreateUser(u *models.User) (*models.User, error) {
	var count int64
	s.db.Model(&gormUser{}).Where("username = ?", u.Username).Count(&count)
	if count > 0 {
		return nil, fmt.Errorf("用户名已存在")
	}

	g := gormUser{
		Username: u.Username,
		Name:     u.Name,
		Password: u.Password,
		Role:     u.Role,
	}
	if err := s.db.Create(&g).Error; err != nil {
		// 走到这里多为数据库故障（重名已被上面 Count 拦下），记日志便于排查
		log.Println("创建用户写入失败(多为数据库故障):", err)
		return nil, err
	}
	return toUser(&g), nil
}

// CreatePost 发一个帖子，返回建好的帖子（含作者）
func (s *Store) CreatePost(content string, author *models.User) *models.Post {
	g := gormPost{
		Content:  content,
		AuthorID: author.ID,
	}
	if err := s.db.Create(&g).Error; err != nil {
		log.Println("发帖写入失败:", err)
		return nil
	}

	// 作者直接用传入的，省一次查库
	p := toPost(&g)
	p.Author = author
	return p
}

// GetPostByID 按帖子ID找帖子（带作者）；找不到返回 nil
func (s *Store) GetPostByID(id int64) *models.Post {
	var g gormPost
	if err := s.db.Preload("Author").First(&g, id).Error; err != nil {
		return nil
	}
	return toPost(&g)
}

// ListPosts 分页取帖子列表，按发布时间倒序
func (s *Store) ListPosts(page, pageSize int) ([]*models.Post, int) {
	var total int64
	s.db.Model(&gormPost{}).Count(&total) // 总数（软删除的被 Gorm 自动排除）

	var list []gormPost
	s.db.Preload("Author").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&list)

	posts := make([]*models.Post, 0, len(list))
	for i := range list {
		posts = append(posts, toPost(&list[i]))
	}
	return posts, int(total)
}

// CreateComment 发评论；帖子不存在返回错误（调用方据此返回 404）
func (s *Store) CreateComment(c *models.Comment) (*models.Comment, error) {
	var post gormPost
	if err := s.db.First(&post, c.PostID).Error; err != nil {
		return nil, fmt.Errorf("帖子不存在")
	}

	g := gormComment{
		PostID:   c.PostID,
		Content:  c.Content,
		AuthorID: c.Author.ID,
	}
	if err := s.db.Create(&g).Error; err != nil {
		log.Println("评论写入失败:", err)
		return nil, err
	}

	// 帖子评论数 +1（原子自增）
	if err := s.db.Model(&gormPost{}).Where("id = ?", c.PostID).
		UpdateColumn("comment_count", gorm.Expr("comment_count + 1")).Error; err != nil {
		log.Println("帖子评论数+1失败:", err)
	}

	cc := toComment(&g)
	cc.Author = c.Author
	return cc, nil
}

// GetCommentsByPostID 取某帖子的所有评论（按创建时间正序）
func (s *Store) GetCommentsByPostID(postID int64) []*models.Comment {
	var list []gormComment
	s.db.Preload("Author").
		Where("post_id = ?", postID).
		Order("created_at ASC").
		Find(&list)

	comments := make([]*models.Comment, 0, len(list))
	for i := range list {
		comments = append(comments, toComment(&list[i]))
	}
	return comments
}

// ToggleLike 点赞开关：没赞→点赞，赞过→取消；帖子不存在返回错误
func (s *Store) ToggleLike(userID, postID int64) (bool, error) {
	var post gormPost
	if err := s.db.First(&post, postID).Error; err != nil {
		return false, fmt.Errorf("帖子不存在")
	}

	var like gormLike
	err := s.db.Where("user_id = ? AND post_id = ?", userID, postID).First(&like).Error
	if err == nil {
		// 赞过 → 取消：删记录，赞数 -1
		if delErr := s.db.Delete(&like).Error; delErr != nil {
			log.Println("取消点赞失败:", delErr)
		}
		s.db.Model(&gormPost{}).Where("id = ?", postID).
			UpdateColumn("like_count", gorm.Expr("like_count - 1"))
		return false, nil
	}

	// 没赞过 → 点赞；复合唯一索引兜底，并发也不会重复
	if createErr := s.db.Create(&gormLike{UserID: userID, PostID: postID}).Error; createErr != nil {
		log.Println("点赞写入失败(可能并发重复或数据异常):", createErr)
	}
	s.db.Model(&gormPost{}).Where("id = ?", postID).
		UpdateColumn("like_count", gorm.Expr("like_count + 1"))
	return true, nil
}

// GetLikeStatus 批量查当前用户对一堆帖子的点赞状态
func (s *Store) GetLikeStatus(userID int64, postIDs []int64) []models.LikeStatus {
	var likes []gormLike
	s.db.Where("user_id = ? AND post_id IN ?", userID, postIDs).Find(&likes)

	likedSet := make(map[int64]bool)
	for _, l := range likes {
		likedSet[l.PostID] = true
	}

	result := make([]models.LikeStatus, 0, len(postIDs))
	for _, pid := range postIDs {
		result = append(result, models.LikeStatus{PostID: pid, Liked: likedSet[pid]})
	}
	return result
}

// DeletePost 删除帖子：帖子软删，评论和点赞级联真删
func (s *Store) DeletePost(postID int64) error {
	var post gormPost
	if err := s.db.First(&post, postID).Error; err != nil {
		return fmt.Errorf("帖子不存在")
	}

	// 软删帖子（给 deleted_at 打标记，可恢复，列表不再显示）
	if err := s.db.Delete(&gormPost{}, postID).Error; err != nil {
		log.Println("软删帖子失败:", err)
	}

	// 评论和点赞表没有软删除字段，所以是真删
	if err := s.db.Where("post_id = ?", postID).Delete(&gormComment{}).Error; err != nil {
		log.Println("级联删除评论失败:", err)
	}
	if err := s.db.Where("post_id = ?", postID).Delete(&gormLike{}).Error; err != nil {
		log.Println("级联删除点赞失败:", err)
	}
	return nil
}
