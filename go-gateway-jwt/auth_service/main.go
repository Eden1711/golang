package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	pb "go-gateway-jwt/proto"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
)

// Khóa bí mật để ký Token (Tuyệt đối không lộ ra ngoài)
var jwtKey = []byte("bi_mat_khong_the_bat_mi")

type server struct {
	pb.UnimplementedAuthServiceServer
}

func (s *server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	// 1. Giả lập check DB
	// Trong thực tế bạn sẽ query Database ở đây
	if req.Username != "admin" || req.Password != "123456" {
		return &pb.LoginResponse{Error: "Sai tài khoản hoặc mật khẩu"}, nil
	}

	// 2. Tạo JWT Token
	// Token này chứa thông tin user và thời gian hết hạn (15 phút)
	claims := jwt.MapClaims{
		"username": req.Username,
		"role":     "admin",
		"exp":      time.Now().Add(15 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return nil, err
	}

	// 3. Trả Token về
	return &pb.LoginResponse{Token: tokenString}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051") // Chạy ở cổng 50051
	if err != nil {
		log.Fatalf("Lỗi mở cổng: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, &server{})

	fmt.Println("🔐 Auth Service (gRPC) đang chạy tại :50051...")
	s.Serve(lis)
}
