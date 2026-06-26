package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Head struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
type Data map[string]any
type Response struct {
	Head Head `json:"head"`
	Data Data `json:"data"`
}

func GinRegisterAndLogin() {
	router := gin.Default()
	router.POST("/register", func(ctx *gin.Context) {
		username := ctx.PostForm("username")
		password := ctx.PostForm("password")
		res := Response{
			Head: Head{
				Code:    http.StatusOK,
				Message: "",
			},
			Data: Data{
				"total": 10,
				"data":  []any{},
			},
		}
		// 校验
		if username == "" || password == "" {
			res.Head.Code = http.StatusBadRequest
			res.Head.Message = "缺失用户名或密码"
			res.Data = Data{}
			ctx.JSON(http.StatusBadRequest, res)
			return
		}
		res.Head.Message = "写入"
		// 写入
		ctx.JSON(http.StatusOK, res)
	})
	router.POST("/login", func(ctx *gin.Context) {
		username := ctx.PostForm("username")
		password := ctx.PostForm("password")
		// 校验
		if username == "admin" && password == "123" {
			ctx.JSON(http.StatusOK, gin.H{"message": "登录成功"})
		} else {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "密码错误"})
		}
	})
	router.Run(":8080")
}
