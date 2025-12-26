package main

import (
	"go-chat/internal/ws"
	"log"
	"net/http"
)

func main() {
	// 1. Tạo Hub
	hub := ws.NewHub()

	// 2. Chạy Hub trên 1 Goroutine riêng (để nó luôn lắng nghe xử lý tin nhắn)
	go hub.Run()

	// 3. Định tuyến
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(hub, w, r)
	})

	log.Println("Server Chat v2 running on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
