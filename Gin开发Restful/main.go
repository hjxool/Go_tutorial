package main

import (
	"net/http"
	// go get . 命令将获取当前目录中代码的依赖项
	"github.com/gin-gonic/gin"
)

// 序列化为 JSON 时字段的名称
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// 结构体切片 模拟的初始数据
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

// go run . 运行当前目录中的代码
func main() {
	router := gin.Default()          // 初始化一个 Gin 路由器
	router.GET("/albums", getAlbums) // GET 函数将 GET HTTP 方法和 /albums 路径与处理器函数关联起来
	router.POST("/albums", postAlbums)
	router.GET("/albums/:id", getAlbumByID) // 冒号 表示 路径参数

	router.Run("localhost:8080") // Run 函数将路由器附加到 http.Server 并启动服务器
}

// 获取专辑列表
func getAlbums(c *gin.Context) {
	// IndentedJSON 带数据格式 JSON 紧凑
	c.IndentedJSON(http.StatusOK, albums)
}

// 添加到专辑列表
func postAlbums(c *gin.Context) {
	var body album
	// BindJSON 将请求正文绑定到 body
	if err := c.BindJSON(&body); err != nil {
		return
	}
	// append 会分配新的底层数组 并返回一个新的切片引用
	// 因此必须用 albums = ... 来接收返回值 否则修改不会生效
	albums = append(albums, body)
	c.IndentedJSON(http.StatusCreated, body) // 201 状态码
}

// 检索特定专辑
func getAlbumByID(c *gin.Context) {
	// 与路由里参数名一一对应 如/users/:uid/orders/:oid 则uid := c.Param("uid") oid := c.Param("oid")
	id := c.Param("id")

	for _, val := range albums {
		if val.ID == id {
			c.IndentedJSON(http.StatusOK, val)
			return
		}
	}
	// gin.H 即 type H map[string]any 别名
	// 实际开发并不用这样的方式一层层嵌套 而是定义结构体
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "没找到专辑"})
}
