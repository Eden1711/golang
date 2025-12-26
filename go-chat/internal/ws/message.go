package ws

// Message đại diện cho gói tin JSON gửi qua lại
type Message struct {
	Type      string `json:"type"`      // "broadcast" (chat all) hoặc "private" (chat riêng)
	Sender    string `json:"sender"`    // Tên người gửi
	Recipient string `json:"recipient"` // Tên người nhận (nếu chat private)
	Content   string `json:"content"`   // Nội dung tin nhắn
}
