package main

import (
	"github.com/gin-gonic/gin"

	"forum/handler"
	"forum/models"
	"forum/response"
	"forum/store"
)

func main() {
	// 请来领班：创建一个 Gin 引擎（默认已带"日志记录"和"崩溃自动恢复"两个贴心服务）
	r := gin.Default()

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
			LikeCount:    12,
			CommentCount: 3,
			CreatedAt:    "2026-07-12T10:00:00+08:00",
		}

		response.Success(c, response.StatusCreated, post)
	})
	// 开食堂的"仓库大厅"：把用户等数据存进去
	myStore := store.New()

	// 创建"前台领班"（鉴权处理器），并把仓库交给他
	authHandler := handler.NewAuthHandler(myStore)

	// 把"注册"这道菜写进菜单：POST /api/v1/auth/register
	r.POST("/api/v1/auth/register", authHandler.Register)

	// 正式营业：守在 8080 号门
	r.Run(":8080")
}
