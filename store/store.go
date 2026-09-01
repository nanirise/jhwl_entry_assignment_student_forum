package store

import (
	"fmt"
	"sync"

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

	// 互斥锁：防止多人同时操作导致数据冲突（类比"同一个人不能同时拿两份"）
	mu sync.Mutex
}

// New 创建并初始化一个存储实例
func New() *Store {
	return &Store{
		users:   make(map[string]*models.User),
		userSeq: 10000, // 起始 ID，第一个用户是 10001
	}
}

// GetUserByUsername 按 username 查找用户；不存在返回 nil
func (s *Store) GetUserByUsername(username string) *models.User {
	// 加锁，避免并发读写冲突
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.users[username]
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
