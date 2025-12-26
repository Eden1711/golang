package main

import (
	"context"
	"log"
	"os"
	"simple-social/internal/api"
	"simple-social/internal/db"
	"simple-social/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dsn := "postgres://myuser:secret@localhost:5432/simple_social?sslmode=disable"

	if envDSN := os.Getenv("DB_SOURCE"); envDSN != "" {
		dsn = envDSN
	}

	// In ra để kiểm tra xem nó đang kết nối đi đâu
	log.Println("Đang kết nối tới DB:", dsn)

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal("Không thể kết nối DB:", err)
	}
	defer pool.Close()

	// 2. Dependency Injection
	store := db.New(pool)                // Repository
	svc := service.NewUserService(store) // Service
	pvc := service.NewPostService(store) // Service
	server := api.NewServer(svc, pvc)    // Server

	// 3. Chạy
	log.Println("Bank Server running on :8080")
	server.Router.Run(":8080")
}
