package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
)

// 必须用大写 对外开放才能被JSON解析器解析
// json标签只影响json里字段名 go文件中访问时还是用大写
type RegisterBody struct {
	Username string `json:"username" validate:"required,min=1"`
	Password string `json:"password" validate:"required,min=6"`
}

func LinkToMysql() {
	validate := validator.New()
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/maindb")
	if err != nil {
		log.Fatal(err)
	}
	// 用defer关键词延迟关闭连接 这样做的好处是不论后续程序是否异常退出 都必定执行defer
	defer db.Close()
	// 测试连接 db.Ping 返回error 如果为nil表示连接正常 否则表示连接失败
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("连接数据库成功")
	router := gin.Default()
	router.POST("/register", func(ctx *gin.Context) {
		// PostForm 是解析 Form-Data 格式的请求体的
		// username := ctx.PostForm("username")
		// password := ctx.PostForm("password")

		// 解析JSON请求体 需要先定义请求体结构 可以引入库来辅助校验
		var body RegisterBody
		// 然后用 BindJSON 把JSON请求体解析到结构体
		if err := ctx.BindJSON(&body); err != nil {
			ctx.JSON(http.StatusOK, Response{
				Head: Head{
					Code:    http.StatusBadRequest,
					Message: "请求体解析失败",
				},
			})
			return
		}
		// 结构体字段值是否缺失等校验
		if err := validate.Struct(body); err != nil {
			ctx.JSON(http.StatusOK, Response{
				Head: Head{
					Code:    http.StatusBadRequest,
					Message: err.Error(),
				},
			})
			return
		}

		// 创建表（如果不存在）表不存在插入数据会报错
		// 字段名 类型 限制
		_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
			id INT AUTO_INCREMENT PRIMARY KEY,
			username VARCHAR(50) NOT NULL,
			password VARCHAR(100) NOT NULL
		)`)
		if err != nil {
			ctx.JSON(http.StatusOK, Response{
				Head: Head{
					Code:    http.StatusBadRequest,
					Message: "创建表错误",
				},
			})
			return
		}
		// 插入数据
		// (username, password) 选择列 (?, ?) 占位符
		_, err = db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", body.Username, body.Password)
		if err != nil {
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
				Message: "数据插入成功",
			},
		})
	})
	// Run 会阻塞程序 因此当前函数不会返回 因此defer db.Close()不会执行 除非程序意外结束
	// 所以在路由的回调函数中 db可以接着用
	router.Run(":8080")
}
