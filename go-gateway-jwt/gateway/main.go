package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	pb "go-gateway-jwt/proto"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var jwtKey = []byte("bi_mat_khong_the_bat_mi") // Phải khớp với Auth Service

func main() {
	// 1. Kết nối tới Auth Service qua gRPC
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Không kết nối được Auth Service: %v", err)
	}
	defer conn.Close()
	authClient := pb.NewAuthServiceClient(conn)

	r := gin.Default()

	// 2. API Public: Đăng nhập (Ai cũng gọi được)
	r.POST("/login", func(c *gin.Context) {
		var loginData struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.BindJSON(&loginData); err != nil {
			c.JSON(400, gin.H{"error": "Dữ liệu sai"})
			return
		}

		// Gọi gRPC sang Auth Service
		resp, err := authClient.Login(context.Background(), &pb.LoginRequest{
			Username: loginData.Username,
			Password: loginData.Password,
		})

		if err != nil {
			c.JSON(500, gin.H{"error": "Lỗi hệ thống"})
			return
		}
		if resp.Error != "" {
			c.JSON(401, gin.H{"error": resp.Error})
			return
		}

		c.JSON(200, gin.H{"token": resp.Token})
	})

	// 3. API Private: Cần có Token mới được vào
	// Sử dụng Middleware AuthMiddleware tự viết ở dưới
	protected := r.Group("/admin")
	protected.Use(AuthMiddleware())
	{
		protected.GET("/dashboard", func(c *gin.Context) {
			// Lấy thông tin user đã lưu trong context ở Middleware
			username := c.MustGet("username").(string)
			c.JSON(200, gin.H{"message": "Chào sếp " + username + "! Đây là dữ liệu mật."})
		})
	}

	fmt.Println("🌐 API Gateway đang chạy tại :8080...")
	r.Run(":8080")
}

// --- MIDDLEWARE KIỂM TRA JWT ---
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Lấy token từ header: "Authorization: Bearer <token>"
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Chưa đăng nhập (Thiếu Header)"})
			c.Abort()
			return
		}

		// Cắt bỏ chữ "Bearer "
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse và kiểm tra Token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"error": "Token hết hạn hoặc không hợp lệ"})
			c.Abort()
			return
		}

		// Lấy thông tin user từ Token (Claims)
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// Lưu username vào context để các hàm sau dùng
			c.Set("username", claims["username"])
		}

		c.Next() // Cho phép đi tiếp
	}
}
