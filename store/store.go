package store

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"forum/models"
)

// Store 简单的内存存储：专门负责"存数据、找数据"
// 现在用内存(memory)保存，后续 Day7 再换成数据库(Gorm)，但接口保持一致
type Store struct {
	// 用 map 存用户：把 username 当"钥匙"，存对应的 User
	// map 是 Go 的"键值对字典"，类比"用学号查人"
	users map[string]*models.User
	// 下一个要分配的 user ID，从 10001 开始（文档示例用户 ID 是这个）
	userSeq int64

	// 帖子货架：把帖子 ID 当"钥匙"，存对应的 Post
	posts map[int64]*models.Post
	// 下一个帖子 ID，从 20001 开始（文档示例帖子 ID 是这个）
	postSeq int64

	// 评论货架：钥匙是帖子 ID，值是"这个帖子的一列评论"
	comments map[int64][]*models.Comment
	// 下一个评论 ID，从 30001 开始（文档示例评论 ID 是这个）
	commentSeq int64

	// 互斥锁：防止多人同时操作导致数据冲突（类比"同一个人不能同时拿两份"）
	mu sync.Mutex
}

// New 创建并初始化一个存储实例
func New() *Store {
	return &Store{
		users:   make(map[string]*models.User),
		userSeq: 10000, // 起始 ID，第一个用户是 10001

		posts:   make(map[int64]*models.Post),
		postSeq: 20000, // 帖子 ID 从 20001 开始

		comments:   make(map[int64][]*models.Comment),
		commentSeq: 30000, // 评论 ID 从 30001 开始
	}
}

// GetUserByUsername 按 username 查找用户；不存在返回 nil
func (s *Store) GetUserByUsername(username string) *models.User {
	// 加锁，避免并发读写冲突
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.users[username]
}

// GetUserByID 按用户 ID 查找用户；不存在返回 nil
// 因为 map 是按 username 存的，所以这里需要遍历找出 ID 匹配的用户
func (s *Store) GetUserByID(id int64) *models.User {
	s.mu.Lock()
	defer s.mu.Unlock()
	// 遍历 map 里所有用户，找到 ID 匹配的那一个
	for _, u := range s.users {
		if u.ID == id {
			return u
		}
	}
	return nil // 没找到
}

// CreateUser 创建一个新用户，返回创建好的用户（含分配好的 ID 和哈希后的密码）
// 如果用户名已存在，返回错误
func (s *Store) CreateUser(u *models.User) (*models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 检查用户名是否已存在
	if _, exists := s.users[u.Username]; exists {
		return nil, fmt.Errorf("用户名已存在")
	}

	// 分配唯一的用户 ID 并自增
	s.userSeq++
	u.ID = s.userSeq

	// 存进字典
	s.users[u.Username] = u
	return u, nil
}

// now 返回当前时间的标准格式字符串，例如 2026-09-03T12:00:00+08:00
// 这个格式叫 RFC3339，前端和时间库都认
func now() string {
	return time.Now().Format(time.RFC3339)
}

// CreatePost 建一个帖子：接收正文 + 作者，自己完成"编号 + 记时间 + 上架"
// 仓库不关心外面发生了什么，只管"给我料，我入库"
func (s *Store) CreatePost(content string, author *models.User) *models.Post {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 发新台号：第一个帖子是 20001
	s.postSeq++
	post := &models.Post{
		ID:        s.postSeq,
		Content:   content,
		Author:    author,
		CreatedAt: now(), // 记下当前时间
	}

	// 按帖子 ID 当钥匙，放上货架
	s.posts[post.ID] = post
	return post
}

// GetPostByID 按帖子 ID 找帖子；找不到返回 nil
func (s *Store) GetPostByID(id int64) *models.Post {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.posts[id]
}

// ListPosts 分页取帖子列表：按发布时间倒序（最新在最上面）
// 返回：当前这一页的帖子 + 帖子总数（总数给前端算共几页用）
func (s *Store) ListPosts(page, pageSize int) ([]*models.Post, int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 把货架上所有帖子装进一个切片
	posts := make([]*models.Post, 0, len(s.posts))
	for _, p := range s.posts {
		posts = append(posts, p)
	}

	// 按时间倒序排。created_at 全是同一格式(RFC3339)，字符串比较就是时间比较
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].CreatedAt > posts[j].CreatedAt
	})

	total := len(posts)

	// 算出这一页的起止下标：起点 = (第几页-1)*每页条数
	start := (page - 1) * pageSize
	if start >= total {
		return []*models.Post{}, total // 超出范围，返回空页
	}
	end := start + pageSize
	if end > total {
		end = total // 最后一页可能不满，收拢到末尾
	}
	return posts[start:end], total
}

// CreateComment 在指定帖子下发一条评论：帖子评论数 +1，并存下这条评论
// 帖子不存在时返回错误（调用方据此返回 404）
func (s *Store) CreateComment(c *models.Comment) (*models.Comment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 先确认帖子在不在，不在直接报错
	post := s.posts[c.PostID]
	if post == nil {
		return nil, fmt.Errorf("帖子不存在")
	}

	// 发新评论号、记时间，塞进该帖子的评论列
	s.commentSeq++
	c.ID = s.commentSeq
	c.CreatedAt = now()
	s.comments[c.PostID] = append(s.comments[c.PostID], c)

	post.CommentCount++ // 帖子的评论数 +1
	return c, nil
}

// GetCommentsByPostID 取某个帖子的所有评论（按发布顺序，即创建时间正序）
// 返回一份副本，避免外部拿到内部切片后并发改出问题
func (s *Store) GetCommentsByPostID(postID int64) []*models.Comment {
	s.mu.Lock()
	defer s.mu.Unlock()

	list := s.comments[postID]
	r := make([]*models.Comment, len(list))
	copy(r, list)
	return r
}
