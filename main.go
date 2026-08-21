package main

import (
	"github.com/gin-gonic/gin"
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

	// 正式营业：守在 8080 号门
	r.Run(":8080")
}
