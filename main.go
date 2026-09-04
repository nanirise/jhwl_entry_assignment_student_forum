package main

import (
	"github.com/gin-gonic/gin"

	"forum/config"
	"forum/handler"
	"forum/middleware"
	"forum/models"
	"forum/response"
	"forum/store"
)

func main() {
	// 创建 Gin 引擎（自带日志与崩溃恢复）
	r := gin.Default()

	// 全局中间件：跨域 CORS（解决网页版 Apifox 的拦截）
	r.Use(middleware.CORS())

	// ---- 演示接口（早期练手用）----
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
	r.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "你好，世界"})
	})
	r.GET("/users/demo", func(c *gin.Context) {
		user := models.User{
			ID:       10001,
			Username: "20260001",
			Name:     "张三",
			Role:     "student",
		}
		response.Success(c, response.StatusCreated, user)
	})
	r.GET("/posts/demo", func(c *gin.Context) {
		post := models.Post{
			ID:           20001,
			Content:      "欢迎报名技术部暑假招新",
			Author:       &models.User{ID: 10001, Username: "20260001", Name: "张三", Role: "student"},
			LikeCount:    12,
			CommentCount: 3,
			CreatedAt:    "2026-07-12T10:00:00+08:00",
		}
		response.Success(c, response.StatusCreated, post)
	})

	// 读配置（数据库地址等，密码在本地私有文件）→ 连库 → 自动建表
	cfg := config.Load()
	db, err := store.Connect(cfg)
	if err != nil {
		panic("连接数据库失败: " + err.Error())
	}
	if err := store.Migrate(db); err != nil {
		panic("建表失败: " + err.Error())
	}

	// 装配：创建仓库和四个处理器，注入仓库（依赖注入）
	myStore := store.New(db)
	authHandler := handler.NewAuthHandler(myStore)
	postHandler := handler.NewPostHandler(myStore)
	likeHandler := handler.NewLikeHandler(myStore)
	adminHandler := handler.NewAdminHandler(myStore)

	// ---- 公开接口（无需登录）----
	r.POST("/api/v1/auth/register", authHandler.Register)
	r.POST("/api/v1/auth/login", authHandler.Login)

	// ---- 受保护接口（走 Auth 鉴权中间件）----
	protected := r.Group("/api/v1", middleware.Auth())
	{
		protected.GET("/auth/me", authHandler.Me)

		protected.POST("/posts", postHandler.Create)
		protected.GET("/posts", postHandler.List)
		protected.GET("/posts/:post_id", postHandler.Get)
		protected.POST("/posts/:post_id/comment", postHandler.CreateComment)

		// ⑥ 删除自己的帖子
		protected.DELETE("/posts/:post_id", postHandler.Delete)

		// ⑧ 点赞/取消
		protected.POST("/posts/:post_id/like", likeHandler.Like)

		// ⑩ 批量点赞状态
		protected.POST("/posts/likes", likeHandler.GetLikeStatus)

		// ⑦ 管理员删除任意帖子
		protected.DELETE("/admin/posts/:post_id", adminHandler.DeletePost)
	}

	// 启动服务，监听 8080
	r.Run(":8080")
}
