package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	soLuongNguoiMua := 1000
	url := "http://localhost:8080/buy"

	var wg sync.WaitGroup
	start := time.Now()

	fmt.Printf("🚀 Bắt đầu thả %d người vào mua vé...\n", soLuongNguoiMua)

	for i := 0; i < soLuongNguoiMua; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			http.Post(url, "application/json", nil)
		}()
	}

	wg.Wait()
	fmt.Printf("✅ Xong! Hết thời gian: %v\n", time.Since(start))
}
