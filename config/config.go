package config

import (
	"github.com/spf13/viper"
)

// Config 集中存放数据库等配置
type Config struct {
	DBHost     string // 数据库地址
	DBPort     string // 端口
	DBUser     string // 用户名
	DBPassword string // 密码
	DBName     string // 库名
}

// Load 读取配置：先读可提交的 config.yaml，再读本地私有 config.local.yaml 覆盖（用来放密码）
func Load() *Config {
	v := viper.New()

	// 给"非敏感"的常见默认值；密码不放默认，放本地私有文件里
	v.SetDefault("db.host", "127.0.0.1")
	v.SetDefault("db.port", "3306")
	v.SetDefault("db.user", "root")
	v.SetDefault("db.name", "forum")

	// 先读可提交的 config.yaml（存在则读，不存在忽略）
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.SetConfigName("config")
	_ = v.ReadInConfig()

	// 再读本地私有 config.local.yaml，用它的值覆盖（密码在这里）
	v.SetConfigName("config.local")
	_ = v.MergeInConfig()

	return &Config{
		DBHost:     v.GetString("db.host"),
		DBPort:     v.GetString("db.port"),
		DBUser:     v.GetString("db.user"),
		DBPassword: v.GetString("db.password"),
		DBName:     v.GetString("db.name"),
	}
}
