package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "go-grpc-demo/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 👇 THAY ĐỔI Ở ĐÂY:
	// Dùng grpc.NewClient thay vì grpc.Dial
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Không thể khởi tạo client: %v", err)
	}
	defer conn.Close()

	// Tạo Client từ connection
	c := pb.NewUserServiceClient(conn)

	// Gọi hàm GetUser
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Thử lấy ID = 1
	r, err := c.GetUser(ctx, &pb.UserRequest{Id: 1})
	if err != nil {
		log.Fatalf("Lỗi gọi hàm: %v", err)
	}

	fmt.Printf("✅ Kết quả Server trả về:\n")
	fmt.Printf("   - Tên: %s\n", r.Name)
	fmt.Printf("   - Email: %s\n", r.Email)
}
