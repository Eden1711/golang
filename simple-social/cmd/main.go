package main

import (
	"context"
	"log"
	"simple-social/internal/api"
	"simple-social/internal/db"
	"simple-social/internal/service"
	"simple-social/util"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Không thể load config:", err)
	}

	log.Println("Đang kết nối tới DB:", config.DBSource)

	pool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("Không thể kết nối DB:", err)
	}
	defer pool.Close()

	// 2. Dependency Injection
	store := db.New(pool)                        // Repository
	svc := service.NewUserService(store, config) // Service
	pvc := service.NewPostService(store)         // Service
	server := api.NewServer(config, svc, pvc)    // Server

	// 3. Chạy
	log.Println("Bank Server running on: ", config.ServerAddress)
	server.Router.Run(config.ServerAddress)
}
