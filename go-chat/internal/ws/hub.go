package ws

import "encoding/json"

type Hub struct {
	// Danh sách các client đang kết nối (Key là con trỏ Client, Value là true)
	// 👇 Sửa: Map từ Username (string) sang Client
	clients map[string]*Client

	// 👇 Sửa: Kênh nhận Message struct thay vì []byte
	broadcast chan Message
	// Kênh đăng ký client mới
	register chan *Client

	// Kênh hủy đăng ký (khi client ngắt kết nối)
	unregister chan *Client
}

func (h *Hub) pushUserList() {
	var users []string
	for name := range h.clients {
		users = append(users, name)
	}

	// Tạo JSON danh sách
	listJSON, _ := json.Marshal(users)

	// Tạo message đặc biệt
	msg := Message{
		Type:    "user_list",
		Content: string(listJSON), // Nhét list vào biến Content
		Sender:  "System",
	}

	// Gửi cho tất cả (Broadcast code cũ)
	bytes, _ := json.Marshal(msg)
	for _, client := range h.clients {
		select {
		case client.send <- bytes:
		default:
			close(client.send)
			delete(h.clients, client.username)
		}
	}
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string]*Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			// Đăng ký: Lưu vào map theo tên
			h.clients[client.username] = client
			h.pushUserList()
		case client := <-h.unregister:
			// Hủy đăng ký: Xóa theo tên
			if _, ok := h.clients[client.username]; ok {
				delete(h.clients, client.username)
				close(client.send)
				h.pushUserList()
			}

		case msg := <-h.broadcast:
			// 🔥 LOGIC ROUTING NẰM Ở ĐÂY 🔥

			// Biến struct thành JSON bytes để gửi đi
			bytes, _ := json.Marshal(msg)

			if msg.Type == "private" {
				// 1. Chat Mật: Chỉ gửi cho người nhận
				if receiver, ok := h.clients[msg.Recipient]; ok {
					select {
					case receiver.send <- bytes:
					default:
						close(receiver.send)
						delete(h.clients, msg.Recipient)
					}
				}
			} else {
				// 2. Chat Public: Gửi cho tất cả (Broadcast)
				for _, client := range h.clients {
					// Đừng gửi lại cho chính người nói (Optional)
					// if client.username == msg.Sender { continue }

					select {
					case client.send <- bytes:
					default:
						close(client.send)
						delete(h.clients, client.username)
					}
				}
			}
		}
	}
}
