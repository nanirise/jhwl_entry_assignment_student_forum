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
	// 请来领班：创建一个 Gin 引擎（默认已带"日志记录"和"崩溃自动恢复"两个贴心服务）
	r := gin.Default()

	// 给所有请求统一加上跨域"允许"头（解决网页版 Apifox 的 CORS 拦截）
	r.Use(middleware.CORS())

	// 定好菜单：当顾客用 GET 方式请求 "/ping" 时，领班带他去这个窗口
	r.GET("/ping", func(c *gin.Context) {
		// c 是领班递给你的"顾客上下文"，装着这次请求的全部信息
		// c.JSON(状态码, 数据)：领班自动把数据包装成 JSON 端给顾客
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/hello", func(c *gin.Context) {
		// c 是领班递给你的"顾客上下文"，装着这次请求的全部信息
		// c.JSON(状态码, 数据)：领班自动把数据包装成 JSON 端给顾客
		c.JSON(200, gin.H{
			"message": "你好，世界",
		})
	})

	// 演示：返回一个"虚构的用户"，看看结构体 + 统一响应长什么样
	r.GET("/users/demo", func(c *gin.Context) {
		// 照 models.User 模板[填一张档案]（明天注册接口才真的存起来）
		user := models.User{
			ID:       10001,
			Username: "20260001",
			Name:     "张三",
			Role:     "student",
		}
		// 用统一纸箱装好，按文档要求"创建成功"用 201
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
	// 加载配置：数据库地址等从 config.yaml / config.local.yaml（密码在后者）读进来
	cfg := config.Load()

	// 连上 MySQL 数据库；连不上直接退出（错误会打在日志里，方便排查）
	db, err := store.Connect(cfg)
	if err != nil {
		panic("连接数据库失败: " + err.Error())
	}

	// 让 Gorm 照着四张"卡片"把表自动建好（表已存在会跳过，重复跑也安全）
	if err := store.Migrate(db); err != nil {
		panic("建表失败: " + err.Error())
	}

	// 开食堂的"仓库大厅"：把连好库的 db 交进去，数据从此存进 MySQL，重启也在
	myStore := store.New(db)

	// 创建"前台领班"（鉴权处理器），并把仓库交给他
	authHandler := handler.NewAuthHandler(myStore)
	// 创建"帖子管家"（帖子/评论处理器）
	postHandler := handler.NewPostHandler(myStore)
	// 创建"点赞管家"（点赞/批量状态处理器）
	likeHandler := handler.NewLikeHandler(myStore)
	// 创建"管理员管家"（管理员专属处理器）
	adminHandler := handler.NewAdminHandler(myStore)

	// 把"注册"这道菜写进菜单：POST /api/v1/auth/register
	r.POST("/api/v1/auth/register", authHandler.Register)

	// 把"登录"这道菜写进菜单：POST /api/v1/auth/login
	r.POST("/api/v1/auth/login", authHandler.Login)

	// 受保护的一组接口：站在这条"走廊"前的是安检门(鉴权中间件)
	// 后面的接口都会先通过 middleware.Auth() 验 JWT 手环
	protected := r.Group("/api/v1", middleware.Auth())
	{
		// 需登录才能访问：返回当前登录用户的信息
		protected.GET("/auth/me", authHandler.Me)

		// 帖子相关：发帖、列表、详情、发评论
		protected.POST("/posts", postHandler.Create)
		protected.GET("/posts", postHandler.List)
		protected.GET("/posts/:post_id", postHandler.Get)
		protected.POST("/posts/:post_id/comment", postHandler.CreateComment)

		// ⑥ 删除自己的帖子
		protected.DELETE("/posts/:post_id", postHandler.Delete)

		// ⑧ 点赞/取消点赞（开关）
		protected.POST("/posts/:post_id/like", likeHandler.Like)

		// ⑩ 批量查询当前用户对多个帖子的点赞状态
		protected.POST("/posts/likes", likeHandler.GetLikeStatus)

		// ⑦ 管理员删除任意帖子
		protected.DELETE("/admin/posts/:post_id", adminHandler.DeletePost)
	}

	// 正式营业：守在 8080 号门
	r.Run(":8080")
}
