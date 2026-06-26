package main

import (
	"fmt"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to Home!")
}
func userHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		fmt.Fprintf(w, "Get User Info")
	} else if r.Method == http.MethodPost {
		fmt.Fprintf(w, "Create User")
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
func HttpServer() {
	http.HandleFunc("/", homeHandler) // 访问根路径时 调用对应函数
	http.HandleFunc("/user", userHandler)

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil) // 启动服务 监听 8080 端口
}
