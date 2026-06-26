package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func LoginMain() {
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/maindb")
	if err != nil {
		fmt.Println("数据库连接失败:", err)
		return
	}
	defer db.Close()
	router := gin.Default()
	router.POST("/login", LoginHandler(db))
	router.Run(":8080")
}

func LoginHandler(db *sql.DB) gin.HandlerFunc {
	validate := validator.New()
	return func(ctx *gin.Context) {
		var body struct {
			Username string `json:"username" validate:"required,min=1"`
			Password string `json:"password" validate:"required,min=6"`
		}
		if err := ctx.BindJSON(&body); err != nil {
			ctx.JSON(http.StatusOK, Response{
				Head: Head{
					Code:    http.StatusBadRequest,
					Message: err.Error(),
				},
			})
			return
		}
		if err := validate.Struct(body); err != nil {
			ctx.JSON(http.StatusOK, Response{
				Head: Head{
					Code:    http.StatusBadRequest,
					Message: "用户名或密码格式错误",
				},
			})
			return
		}
		ok, err := LoginServer(db, body.Username, body.Password)
		if err != nil {
			ctx.JSON(http.StatusOK, Response{
				Head: Head{
					Code:    http.StatusBadRequest,
					Message: "查询用户错误",
				},
			})
			return
		}
		if !ok {
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
				Message: "登录成功",
			},
		})
	}
}

func LoginServer(db *sql.DB, username, password string) (bool, error) {
	user, err := LoginRepo(db, username)
	if err != nil {
		return false, err
	}
	return user.Password == password, nil
}

type LoginUser struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginRepo(db *sql.DB, username string) (LoginUser, error) {
	var user LoginUser
	// QueryRow 单行查询
	row := db.QueryRow("SELECT id, username, password FROM users WHERE username = ?", username)
	// Scan 按 SELECT顺序赋值
	err := row.Scan(&user.ID, &user.Username, &user.Password)
	return user, err
}
