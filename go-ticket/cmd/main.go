package main

import (
	"context"
	"fmt"
	"go-ticket/internal/db"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

const dbSource = "postgres://myuser:secret@localhost:5432/go_ticket?sslmode=disable"

func main() {
	pool, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()
	queries := db.New(pool)

	r := gin.Default()
	// 4. API Mua vé (Chứa logic lỗi Race Condition)
	r.POST("/buy", func(c *gin.Context) {
		ctx := c.Request.Context()

		// 🔥 1. BẮT ĐẦU TRANSACTION
		tx, err := pool.Begin(ctx)
		if err != nil {
			c.JSON(500, gin.H{"error": "Lỗi tạo Transaction"})
			return
		}
		// Luôn Rollback nếu có lỗi xảy ra giữa đường (để giải phóng khóa)
		defer tx.Rollback(ctx)

		// 🔥 2. Tạo query gắn với Transaction này (Quan trọng!)
		// Nếu dùng `queries.Get...` thường thì nó sẽ chạy ngoài transaction -> Không lock được
		qtx := queries.WithTx(tx)

		// 🔥 3. Lấy vé và KHÓA (SELECT FOR UPDATE)
		// Các request khác đến đây sẽ bị treo (đứng chờ)
		ticket, err := qtx.GetTicketForUpdate(ctx, 1)
		if err != nil {
			c.JSON(500, gin.H{"error": "Lỗi DB"})
			return
		}

		// 4. Kiểm tra vé
		if ticket.Quantity > 0 {
			// Ngủ vẫn ok, vì ta đã khóa DB rồi, thằng khác phải chờ ta ngủ dậy
			time.Sleep(10 * time.Millisecond)

			// 5. Trừ vé
			newQuantity := ticket.Quantity - 1
			err = qtx.UpdateTicket(ctx, db.UpdateTicketParams{
				ID:       1,
				Quantity: newQuantity,
			})
			if err != nil {
				c.JSON(500, gin.H{"error": "Update thất bại"})
				return
			}

			// 🔥 6. COMMIT TRANSACTION (Xác nhận và Mở khóa cho người sau)
			err = tx.Commit(ctx)
			if err != nil {
				c.JSON(500, gin.H{"error": "Commit thất bại"})
				return
			}

			c.JSON(200, gin.H{"msg": "Mua thành công", "con_lai": newQuantity})
		} else {
			// Hết vé cũng phải Commit để nhả khóa (hoặc để defer Rollback lo)
			c.JSON(400, gin.H{"msg": "Hết vé rồi!"})
		}
	})

	// 5. API Reset lại vé về 100 (để test lại cho nhanh)
	r.POST("/reset", func(c *gin.Context) {
		queries.UpdateTicket(c.Request.Context(), db.UpdateTicketParams{ID: 1, Quantity: 100})
		c.JSON(200, gin.H{"msg": "Đã reset về 100 vé"})
	})

	r.GET("/ticket", func(c *gin.Context) {
		listTicket, err := queries.GetAllTicket(context.Background())
		if err != nil {
			log.Println(err)
			c.JSON(500, gin.H{"error": "Lỗi DB"})
			return
		}
		c.JSON(200, listTicket)
	})

	fmt.Println("Server đang chạy tại cổng 8080...")
	r.Run(":8080")
}
