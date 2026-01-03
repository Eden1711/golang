package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"go-gateway-jwt/pkg/telemetry"
	pb "go-gateway-jwt/proto"

	"github.com/golang-jwt/jwt/v5"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
)

// Khóa bí mật để ký Token (Tuyệt đối không lộ ra ngoài)
var jwtKey = []byte("bi_mat_khong_the_bat_mi")

type server struct {
	pb.UnimplementedAuthServiceServer
}

func (s *server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	// 1. TẠO SPAN CON (ĐỂ SOI CHI TIẾT)
	// Giả sử ta muốn đo xem việc "Check DB" tốn bao lâu
	tracer := otel.Tracer("auth-service")
	// Tạo 1 đoạn trace con tên là "database_check"
	ctx, span := tracer.Start(ctx, "database_check")

	// Giả vờ ngủ 500ms để mô phỏng DB chậm
	time.Sleep(500 * time.Millisecond)

	// Kết thúc đo đạc
	span.End()

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
	// KHỞI TẠO TRACER
	shutdown := telemetry.InitTracer("auth-service", "jaeger:4317")
	defer shutdown(context.Background())

	lis, err := net.Listen("tcp", ":50051") // Chạy ở cổng 50051
	if err != nil {
		log.Fatalf("Lỗi mở cổng: %v", err)
	}

	// GẮN INTERCEPTOR CHO SERVER GRPC
	// Để nó hiểu và nối tiếp TraceID từ Gateway gửi sang
	s := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()), // 👈 QUAN TRỌNG
	)
	pb.RegisterAuthServiceServer(s, &server{})

	fmt.Println("🔐 Auth Service (gRPC) đang chạy tại :50051...")
	s.Serve(lis)
}
