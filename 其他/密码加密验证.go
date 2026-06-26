package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func EncryptPassword() {
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/maindb")
	if err != nil {
		fmt.Println("数据库连接错误")
		return
	}
	defer db.Close()

	r := gin.Default()
	r.POST("/register", func(ctx *gin.Context) {
		var body RegisterBody
		if err := ctx.BindJSON(&body); err != nil {
			ctx.JSON(http.StatusOK, Response{
				Head: Head{
					Code:    http.StatusBadRequest,
					Message: "请求体解析失败",
				},
			})
			return
		}
		// 加密密码
		bytes, err1 := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
		password := string(bytes)
		if err1 != nil {
			ctx.JSON(http.StatusOK, Response{
				Head: Head{
					Code:    http.StatusInternalServerError,
					Message: "密码加密失败",
				},
			})
			return
		}
		// 保存加密密码
		_, err2 := db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", body.Username, password)
		if err2 != nil {
			ctx.JSON(http.StatusOK, Response{
				Head: Head{
					Code:    http.StatusBadRequest,
					Message: "插入数据错误",
				},
			})
			return
		}
		ctx.JSON(http.StatusOK, Response{
			Head: Head{
				Code:    http.StatusOK,
				Message: "注册成功",
			},
		})
	})
	r.POST("/verify", func(ctx *gin.Context) {
		var body RegisterBody
		if err := ctx.BindJSON(&body); err != nil {
			ctx.JSON(http.StatusOK, Response{
				Head: Head{
					Code:    http.StatusBadRequest,
					Message: "请求体解析失败",
				},
			})
			return
		}
		// 查询用户名对应的加密密码
		var passwordHash string
		err := db.QueryRow("SELECT password FROM users WHERE username = ?", body.Username).Scan(&passwordHash)
		if err != nil {
			ctx.JSON(http.StatusOK, Response{
				Head: Head{
					Code:    http.StatusBadRequest,
					Message: "查询用户错误",
				},
			})
			return
		}
		// 验证密码
		err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(body.Password))
		if err != nil {
			ctx.JSON(http.StatusOK, Response{
				Head: Head{
					Code:    http.StatusBadRequest,
					Message: "密码错误",
				},
			})
			return
		}
		ctx.JSON(http.StatusOK, Response{
			Head: Head{
				Code:    http.StatusOK,
				Message: "验证成功",
			},
		})
	})

	r.Run(":8080")
}
