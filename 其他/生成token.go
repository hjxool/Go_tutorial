package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken() {
	r := gin.Default()

	secretKey := []byte("sercet_key")
	getToken := func(user LoginUser) (string, error) {
		claims := jwt.MapClaims{
			"username": user.Username,
			"password": user.Password,
			"exp":      time.Now().Add(time.Second * 20).Unix(), // 20秒后过期
		}
		// 还没有签名的Token对象 记录了签名算法和声明
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		// 签名并获取完整的token字符串 对称加密算法要求传入[]byte参数作为密钥
		return token.SignedString(secretKey)
	}
	loginHandler := func(ctx *gin.Context) {
		var body LoginUser
		if err := ctx.BindJSON(&body); err != nil {
			ctx.JSON(http.StatusOK, Response{
				Head: Head{
					Code:    http.StatusBadRequest,
					Message: err.Error(),
				},
			})
			return
		}
		token, err := getToken(body)
		if err != nil {
			ctx.JSON(http.StatusOK, Response{
				Head: Head{
					Code:    http.StatusBadRequest,
					Message: err.Error(),
				},
			})
			return
		}
		ctx.JSON(http.StatusOK, Response{
			Head: Head{
				Code:    http.StatusOK,
				Message: "success",
			},
			Data: Data{ // type Data map[string]any
				// 这里要用小写 这样在返回给前端时 json里才是小写 不需要遵守导出规则大写
				"token": token,
			},
		})
	}
	verifyToken := func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		// TrimPrefix会删除前缀 如果没有对应前缀则返回原字符串
		token = strings.TrimPrefix(token, "Bearer ")
		if token == "" {
			ctx.JSON(http.StatusOK, Response{
				Head: Head{
					Code:    http.StatusBadRequest,
					Message: "token不能为空",
				},
			})
			ctx.Abort()
			return
		}
		claims := jwt.MapClaims{} // 声明空结构体
		// 验证token 将解密后的参数放到 claims 里
		_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
			return secretKey, nil
		})
		if err != nil {
			ctx.JSON(http.StatusOK, Response{
				Head: Head{
					Code:    http.StatusBadRequest,
					Message: "无效token",
				},
			})
			ctx.Abort()
			return
		}
		// 把解密出的信息放到 当前请求链上的 *gin.Context中
		ctx.Set("userInfo", map[string]any{
			"username": claims["username"],
			"password": claims["password"],
		})
		ctx.Next()
	}

	r.POST("/login", loginHandler)
	r.GET("/verify", verifyToken, func(ctx *gin.Context) {
		userInfo, exist := ctx.Get("userInfo")
		if !exist {
			ctx.JSON(http.StatusOK, Response{
				Head: Head{
					Code:    http.StatusBadRequest,
					Message: "用户信息不存在",
				},
			})
			return
		}
		ctx.JSON(http.StatusOK, Response{
			Head: Head{
				Code:    http.StatusOK,
				Message: "success",
			},
			Data: Data{
				"username": userInfo.(map[string]any)["username"],
				"password": userInfo.(map[string]any)["password"],
			},
		})
	})
	r.Run(":8080")
}
