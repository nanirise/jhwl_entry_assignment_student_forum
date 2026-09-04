package store

import (
	"fmt"

	"forum/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Connect 用配置连上 MySQL，返回 *gorm.DB；连不上则返回错误
func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Migrate 照着 gorm_models.go 自动建表；表已存在则跳过，重复跑安全
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&gormUser{}, &gormPost{}, &gormComment{}, &gormLike{})
}
