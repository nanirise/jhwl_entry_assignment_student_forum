package store

import (
	"fmt"

	"forum/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Connect 用配置连上 MySQL，返回一个能查询的通行证 *gorm.DB。
// 连不上就返回错误，绝不让程序"带病跑"。
func Connect(cfg *config.Config) (*gorm.DB, error) {
	// dsn = 连接串：用户名:密码@tcp(地址:端口)/库名?后缀参数
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	// gorm.Open 用 mysql 驱动真正打开连接；db 是通行证，err 是错误
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Migrate 照着 gorm_models.go 里四张"专属卡片"自动建表。
// 表若已存在就跳过，重复跑也安全。
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&gormUser{}, &gormPost{}, &gormComment{}, &gormLike{})
}
