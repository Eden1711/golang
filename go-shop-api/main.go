package main

import (
	"context"
	"go-shop-api/api"
	"go-shop-api/db"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dbSource := "postgresql://postgres:secret@localhost:5432/product_db?sslmode=disable"
	connPool, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal("Cannot connect to db:", err)
	}
	defer connPool.Close()

	// 2. Tạo Store (Dependency)
	store := db.New(connPool)

	// 3. Khởi tạo Server và INJECT Store vào (DI ở đây)
	server := api.NewServer(store)

	// 4. Start Server
	log.Println("Server running on :8080")
	err = server.Start(":8080")
	if err != nil {
		log.Fatal("Cannot start server:", err)
	}
}
