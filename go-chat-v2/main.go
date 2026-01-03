package main

import (
	"flag" // Import thư viện đọc cờ lệnh
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func main() {
	// 1. Đọc tham số cổng từ dòng lệnh (Mặc định là 8080)
	port := flag.String("port", "8080", "Cổng chạy server")
	flag.Parse()

	// 2. Khởi tạo Redis (Biến rdb khai báo bên client.go)
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	hub := newHub()
	go hub.run()

	// 🔴 3. Chạy thêm Goroutine lắng nghe Redis
	go hub.subscribeToRedis()

	r := gin.Default()
	r.GET("/ws", func(c *gin.Context) {
		serveWs(hub, c)
	})

	r.LoadHTMLFiles("index.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	// Chạy trên cổng động
	addr := ":" + *port
	fmt.Printf("Server đang chạy tại http://localhost%s\n", addr)
	r.Run(addr)
}
