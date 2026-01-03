package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

// Khai báo biến Redis toàn cục (hoặc truyền vào struct cũng được)
var rdb *redis.Client

const (
	// Thời gian chờ để ghi tin nhắn ra
	writeWait = 10 * time.Second
	// Thời gian tối đa để nhận pong từ client (check heartbeat)
	pongWait = 60 * time.Second
	// Chu kỳ gửi ping
	pingPeriod = (pongWait * 9) / 10
)

// Cấu hình nâng cấp từ HTTP -> WebSocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Cho phép mọi nguồn (CORS) để dễ test local
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Client struct {
	hub *Hub
	// Kết nối socket thực sự
	conn *websocket.Conn
	// Kênh đệm để giữ tin nhắn cần gửi cho user này
	send chan []byte
}

// 1. Đọc tin từ Browser gửi lên
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	// Cấu hình giới hạn thời gian đọc
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		// 🔴 THAY ĐỔI: Thay vì gửi vào hub.broadcast, ta bắn lên Redis
		// c.hub.broadcast <- message  <-- Cũ (Xóa hoặc comment dòng này)

		// Mới: Publish vào kênh "chat_room"
		err = rdb.Publish(context.Background(), "chat_room", message).Err()
		if err != nil {
			log.Println("Lỗi Redis Publish:", err)
		}
	}
}

// 2. Ghi tin từ Server xuống Browser
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub đã đóng kênh
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// Gửi tin nhắn text
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Xả bộ đệm
			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			// Gửi Ping định kỳ để giữ kết nối không bị đứt
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// Hàm API để nâng cấp kết nối HTTP thành WebSocket
func serveWs(hub *Hub, c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Tạo client mới
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}

	// Đăng ký với Hub
	client.hub.register <- client

	// Chạy 2 goroutine để đọc và ghi song song
	go client.writePump()
	go client.readPump()
}
