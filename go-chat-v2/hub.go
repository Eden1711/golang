package main

import (
	"context"
	"fmt"
)

// Hub quản lý trạng thái của các phòng chat
type Hub struct {
	// Danh sách các Client đang kết nối (Key là con trỏ Client, Value là bool)
	clients map[*Client]bool

	// Kênh nhận yêu cầu đăng ký
	register chan *Client

	// Kênh nhận yêu cầu hủy đăng ký
	unregister chan *Client

	// Kênh phát tin nhắn chung (Broadcast)
	broadcast chan []byte
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// Thêm hàm này vào Hub
func (h *Hub) subscribeToRedis() {
	// Đăng ký nghe kênh "chat_room"
	pubsub := rdb.Subscribe(context.Background(), "chat_room")

	// Lấy channel Go từ Redis PubSub
	ch := pubsub.Channel()

	// Vòng lặp nhận tin từ Redis
	for msg := range ch {
		// Khi Redis báo có tin nhắn mới -> Đẩy vào broadcast nội bộ
		h.broadcast <- []byte(msg.Payload)
	}
}

// Vòng lặp chính của Hub (Chạy trên 1 Goroutine riêng)
func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			// Có người vào -> Lưu vào map
			h.clients[client] = true
			fmt.Println("➕ Một user vừa vào phòng")

		case client := <-h.unregister:
			// Có người ra -> Xóa khỏi map và đóng kênh
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				fmt.Println("➖ Một user đã thoát")
			}

		case message := <-h.broadcast:
			// Có tin nhắn -> Gửi cho TẤT CẢ client
			for client := range h.clients {
				select {
				case client.send <- message:
					// Gửi thành công vào kênh riêng của client đó
				default:
					// Nếu client bị lag/đơ, xóa luôn để đỡ tốn bộ nhớ
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
