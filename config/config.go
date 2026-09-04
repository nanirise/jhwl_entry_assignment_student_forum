package config

import (
	"github.com/spf13/viper"
)

// Config 数据库等配置
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

// Load 读取配置：先读可提交的 config.yaml，再用本地 config.local.yaml 覆盖（密码放后者）
func Load() *Config {
	v := viper.New()

	// 非敏感默认值；密码不放默认，只放在本地私有文件
	v.SetDefault("db.host", "127.0.0.1")
	v.SetDefault("db.port", "3306")
	v.SetDefault("db.user", "root")
	v.SetDefault("db.name", "forum")

	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.SetConfigName("config")
	_ = v.ReadInConfig()

	// 用本地私有文件覆盖（密码在这里，已被 .gitignore 忽略）
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
