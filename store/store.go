package store

import (
	"fmt"
	"log"

	"forum/models"
	"gorm.io/gorm"
)

// Store 仓储：这次换成真·MySQL 当"仓库"，用 Gorm 查询。
// 对外方法的名字、参数、返回值都不变（"换心不换脸"），所以 handler 一个字都不用改。
type Store struct {
	db *gorm.DB // 通行证：所有查库动作都叫它去干
}

// New 创建仓储：传入连好库的 *gorm.DB（在 main 里先 store.Connect 再传进来）
func New(db *gorm.DB) *Store {
	return &Store{db: db}
}

// GetUserByUsername 按学号查用户；查不到返回 nil
func (s *Store) GetUserByUsername(username string) *models.User {
	var g gormUser
	// First = 取"第一条"。没查到会返回 err，我们把它当"不存在"处理，返回 nil
	if err := s.db.Where("username = ?", username).First(&g).Error; err != nil {
		return nil
	}
	return toUser(&g)
}

// GetUserByID 按用户 ID 查用户；查不到返回 nil
func (s *Store) GetUserByID(id int64) *models.User {
	var g gormUser
	if err := s.db.First(&g, id).Error; err != nil {
		return nil
	}
	return toUser(&g)
}

// CreateUser 创建新用户；用户名已存在返回错误（调用方据此返回 409）
func (s *Store) CreateUser(u *models.User) (*models.User, error) {
	// 1. 先数一下库里有没有同名的：Count 统计条数
	var count int64
	s.db.Model(&gormUser{}).Where("username = ?", u.Username).Count(&count)
	if count > 0 {
		return nil, fmt.Errorf("用户名已存在")
	}

	// 2. 组装数据库卡：注意密码(Password)在 handler 里已经哈希过了，这里直接存
	g := gormUser{
		Username: u.Username,
		Name:     u.Name,
		Password: u.Password,
		Role:     u.Role,
	}

	// 3. Create = 插入一行；Gorm 自动把 ID、CreatedAt、UpdatedAt 填进 g
	if err := s.db.Create(&g).Error; err != nil {
		// 走到这里往往不是"重名"(重名已被上面的 Count 拦下)，多是数据库故障；记下真实原因
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

	// 作者直接用传进来的（我们知道发帖人是谁），省一次查库
	p := toPost(&g)
	p.Author = author
	return p
}

// GetPostByID 按帖子 ID 找帖子（带上作者）；找不到返回 nil
func (s *Store) GetPostByID(id int64) *models.Post {
	var g gormPost
	// Preload("Author") = 把"作者"这张关联卡一起查出，塞进 g.Author
	if err := s.db.Preload("Author").First(&g, id).Error; err != nil {
		return nil
	}
	return toPost(&g)
}

// ListPosts 分页取帖子列表，按发布时间倒序（最新在最上面）
func (s *Store) ListPosts(page, pageSize int) ([]*models.Post, int) {
	// 1. 总数：Count 数"还剩多少帖"（软删除的会被 Gorm 自动排除）
	var total int64
	s.db.Model(&gormPost{}).Count(&total)

	// 2. 取这一页：倒序 + 跳过前面若干条(offset) + 只要 pageSize 条(limit)
	var list []gormPost
	s.db.Preload("Author").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&list)

	// 3. 把一堆"数据库卡"翻译成"商品卡"装进切片
	posts := make([]*models.Post, 0, len(list))
	for i := range list {
		posts = append(posts, toPost(&list[i]))
	}
	return posts, int(total)
}

// CreateComment 在指定帖子下发评论；帖子不存在返回错误（调用方据此返回 404）
func (s *Store) CreateComment(c *models.Comment) (*models.Comment, error) {
	// 1. 确认帖子在不在
	var post gormPost
	if err := s.db.First(&post, c.PostID).Error; err != nil {
		return nil, fmt.Errorf("帖子不存在")
	}

	// 2. 插入评论
	g := gormComment{
		PostID:   c.PostID,
		Content:  c.Content,
		AuthorID: c.Author.ID,
	}
	if err := s.db.Create(&g).Error; err != nil {
		log.Println("评论写入失败:", err)
		return nil, err
	}

	// 3. 帖子的评论数 +1（UpdateColumn 直接改数据库那一列）
	if err := s.db.Model(&gormPost{}).Where("id = ?", c.PostID).
		UpdateColumn("comment_count", gorm.Expr("comment_count + 1")).Error; err != nil {
		log.Println("帖子评论数+1失败:", err)
	}

	// 4. 回填作者，返回建好的评论
	cc := toComment(&g)
	cc.Author = c.Author
	return cc, nil
}

// GetCommentsByPostID 取某个帖子的所有评论（按发布顺序，即创建时间正序）
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

// ToggleLike 点赞开关：原本没赞 → 点赞，原本赞了 → 取消。
// 帖子不存在返回错误（调用方据此返回 404）。
func (s *Store) ToggleLike(userID, postID int64) (bool, error) {
	// 1. 帖子得存在，才谈得上赞它
	var post gormPost
	if err := s.db.First(&post, postID).Error; err != nil {
		return false, fmt.Errorf("帖子不存在")
	}

	// 2. 查这个用户赞没赞过这个帖
	var like gormLike
	err := s.db.Where("user_id = ? AND post_id = ?", userID, postID).First(&like).Error
	if err == nil {
		// 赞过 → 取消：删掉那条点赞记录，赞数 -1
		if delErr := s.db.Delete(&like).Error; delErr != nil {
			log.Println("取消点赞失败:", delErr)
		}
		s.db.Model(&gormPost{}).Where("id = ?", postID).
			UpdateColumn("like_count", gorm.Expr("like_count - 1"))
		return false, nil
	}

	// 3. 没赞过 → 点赞：新增记录，赞数 +1
	//    (user_id, post_id) 复合唯一索引，数据库层面保证同一个用户对同一帖不会重复点赞
	if createErr := s.db.Create(&gormLike{UserID: userID, PostID: postID}).Error; createErr != nil {
		log.Println("点赞写入失败(可能并发重复或数据异常):", createErr)
	}
	s.db.Model(&gormPost{}).Where("id = ?", postID).
		UpdateColumn("like_count", gorm.Expr("like_count + 1"))
	return true, nil
}

// GetLikeStatus 批量查"当前用户对一堆帖子"的点赞状态；没赞或不存在都记 false
func (s *Store) GetLikeStatus(userID int64, postIDs []int64) []models.LikeStatus {
	// 一次性把这批帖里"userID 赞过的"都捞出来
	var likes []gormLike
	s.db.Where("user_id = ? AND post_id IN ?", userID, postIDs).Find(&likes)

	// 把"赞过的帖子ID"放进一个集合，方便 O(1) 判断
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

// DeletePost 删除帖子：帖子(软删) + 级联硬删它的评论和点赞记录
func (s *Store) DeletePost(postID int64) error {
	var post gormPost
	if err := s.db.First(&post, postID).Error; err != nil {
		return fmt.Errorf("帖子不存在")
	}

	// 软删帖子：给已删除的行打上 deleted_at 标记（不清数据、可恢复，列表不再显示）
	if err := s.db.Delete(&gormPost{}, postID).Error; err != nil {
		log.Println("软删帖子失败:", err)
	}

	// 级联：评论和点赞表没有 deleted_at 字段，所以对它们是"真删"
	if err := s.db.Where("post_id = ?", postID).Delete(&gormComment{}).Error; err != nil {
		log.Println("级联删除评论失败:", err)
	}
	if err := s.db.Where("post_id = ?", postID).Delete(&gormLike{}).Error; err != nil {
		log.Println("级联删除点赞失败:", err)
	}
	return nil
}
