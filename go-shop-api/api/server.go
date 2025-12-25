package api

import (
	"go-shop-api/db"

	"github.com/gin-gonic/gin"
)

// Server phục vụ HTTP request
// Nó nắm giữ các dependencies cần thiết (ở đây là db.Queries và gin.Engine)
type Server struct {
	store  *db.Queries // Dependency được inject vào
	router *gin.Engine
}

// NewServer: Đây là Constructor
// Ai muốn tạo Server thì PHẢI đưa store vào đây -> Dependency Injection
func NewServer(store *db.Queries) *Server {
	server := &Server{
		store: store,
	}

	router := gin.Default()

	// Đăng ký routes
	// Lưu ý: Các hàm handler giờ là Method của Server (server.createProduct)
	router.POST("/products", server.createProduct)
	router.GET("/products", server.listProducts)
	router.GET("/products/:id", server.getProduct)
	router.PUT("/products/:id", server.updateProduct)
	router.DELETE("/products/:id", server.deleteProduct)
	server.router = router
	return server
}

// Hàm để chạy server (Start)
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
