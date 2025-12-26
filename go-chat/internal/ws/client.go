package ws

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// Cấu hình WebSocket (như cũ)
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Struct đại diện cho 1 người dùng
type Client struct {
	hub      *Hub            // Tham chiếu đến Hub
	conn     *websocket.Conn // Kết nối socket
	send     chan []byte     // Kênh đệm để chứa tin nhắn cần gửi cho user này
	username string
}

// 1. Goroutine ĐỌC: Nhận tin từ Browser -> Đẩy vào Hub
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		// 2. Thay đổi cách đọc: Đọc JSON vào struct Message
		var msg Message
		err := c.conn.ReadJSON(&msg) // Dùng ReadJSON thay cho ReadMessage
		if err != nil {
			break
		}

		// Gán danh tính người gửi (Server tự điền để tránh giả mạo)
		msg.Sender = c.username

		// Đẩy struct Message vào Hub (thay vì []byte)
		c.hub.broadcast <- msg
	}
}

// 2. Goroutine GHI: Nhận tin từ Hub -> Gửi xuống Browser
func (c *Client) writePump() {
	defer c.conn.Close()

	// 👇 SỬA: Dùng for range thay vì for { select {} }
	// Vòng lặp này sẽ chạy liên tục mỗi khi có tin nhắn vào c.send
	// Nó tự động thoát khi kênh c.send bị đóng (close)
	for message := range c.send {
		w, err := c.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			return
		}
		w.Write(message)

		// Gửi các tin nhắn còn tồn đọng trong hàng đợi (nếu có)
		n := len(c.send)
		for i := 0; i < n; i++ {
			w.Write(<-c.send)
		}

		if err := w.Close(); err != nil {
			return
		}
	}

	// Khi vòng lặp kết thúc (nghĩa là kênh c.send đã bị đóng bởi Hub),
	// ta gửi tin nhắn đóng kết nối cho Client biết
	c.conn.WriteMessage(websocket.CloseMessage, []byte{})
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// 3. Lấy tên từ URL: ws://localhost:8080/ws?name=Batman
	username := r.URL.Query().Get("name")
	if username == "" {
		http.Error(w, "Missing 'name' param", 400)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := &Client{
		hub:      hub,
		conn:     conn,
		send:     make(chan []byte, 256),
		username: username, // 👇 Gán tên
	}

	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}
