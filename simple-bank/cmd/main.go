package main

import (
	"context"
	"log"
	"simple-bank/internal/api"
	"simple-bank/internal/db"
	"simple-bank/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// 1. Kết nối DB
	// Nhớ đổi tên DB thành simple_bank
	dsn := "postgres://myuser:secret@localhost:5432/simple_bank?sslmode=disable"
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// 2. Dependency Injection
	store := db.New(pool)                   // Repository
	svc := service.NewAccountService(store) // Service
	server := api.NewServer(svc)            // Server

	// 3. Chạy
	log.Println("Bank Server running on :8080")
	server.Router.Run(":8080")
}
