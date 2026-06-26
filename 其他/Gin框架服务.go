package main

import "github.com/gin-gonic/gin"

func GinServerTest() {
	router := gin.Default() // 创建一个带 Logger 和 Recovery 中间件的路由引擎
	router.GET("/", func(c *gin.Context) {
		c.String(200, "Welcome to Home!")
	})
	router.GET("/user", func(c *gin.Context) {
		c.String(200, "Get User Info")
	})
	router.POST("/user", func(c *gin.Context) {
		c.String(200, "Create User")
	})
	router.Run(":8080") // 启动服务
}
