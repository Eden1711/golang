package main

import (
	"context"
	"fmt"
	pb "go-grpc-demo/proto"
	"log"
	"net"

	"google.golang.org/grpc"
)

// 1. Định nghĩa struct Server (phải khớp với interface sinh ra)
type server struct {
	pb.UnimplementedUserServiceServer
}

// 2. Viết logic cho hàm GetUser
func (s *server) GetUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	fmt.Printf("🔥 Nhận request lấy ID: %d\n", req.Id)

	// Giả lập lấy từ DB (Hardcode cho nhanh)
	if req.Id == 1 {
		return &pb.UserResponse{
			Id:    1,
			Name:  "Batman",
			Email: "batman@gotham.com",
		}, nil
	}

	return &pb.UserResponse{
		Id:    req.Id,
		Name:  "Người lạ",
		Email: "unknown@example.com",
	}, nil
}
func main() {
	// 3. Mở cổng TCP (gRPC chạy trên HTTP/2 nhưng cần TCP listener)
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Không thể mở cổng: %v", err)
	}

	// 4. Khởi tạo gRPC Server
	s := grpc.NewServer()

	// 5. Đăng ký service của mình lên server
	pb.RegisterUserServiceServer(s, &server{})

	fmt.Println("🚀 gRPC Server đang chạy tại :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Lỗi server: %v", err)
	}
}
