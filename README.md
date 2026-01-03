# 🚀 Go High-Performance & Distributed Systems

![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-Enabled-2496ED?logo=docker&logoColor=white)
![Redis](https://img.shields.io/badge/Redis-Pub%2FSub-DC382D?logo=redis&logoColor=white)
![gRPC](https://img.shields.io/badge/gRPC-Protobuf-4285F4?logo=google&logoColor=white)

> Một tập hợp các dự án thực tế mô phỏng các hệ thống Backend chịu tải cao, kiến trúc Microservices và xử lý thời gian thực.
> Dự án này tập trung giải quyết các bài toán khó về: **Concurrency, Race Conditions, Distributed Locking, và Real-time Communication.**

---

## 🏗 Kiến trúc Tổng quan (Architecture)

_(Bạn hãy vẽ một sơ đồ nối các service lại và dán ảnh vào đây. Dùng Excalidraw.com rất đẹp)_

Hệ thống bao gồm các module chính:

| Module                   | Công nghệ & Kỹ thuật chính                    | Bài toán giải quyết                                                                                |
| :----------------------- | :-------------------------------------------- | :------------------------------------------------------------------------------------------------- |
| **Go-Ticket**            | Postgres Lock, **Redis Lua Script**           | Xử lý **Flash Sale** (1 triệu req/s), chặn **Overselling** (Bán lố) khi hàng nghìn người cùng mua. |
| **Go-gateway-jwt**       | **gRPC**, Microservices, **JWT**, API Gateway | Bảo mật hệ thống phân tán, giao tiếp nội bộ siêu tốc, che giấu service sau Gateway.                |
| **Go-Chat / Go-Chat-V2** | **WebSockets**, **Redis Pub/Sub**             | Hệ thống chat Real-time, Scale-out nhiều server (User A ở Server 1 chat với User B ở Server 2).    |

---

## 🔥 Chi tiết Kỹ thuật & Giải pháp (Engineering Decisions)

### 1. Xử lý Race Condition trong Flash Sale

- **Vấn đề:** Khi dùng code thông thường, 1000 request cùng đọc Database thấy `quantity=1` và cùng trừ, dẫn đến bán lố vé.
- **Giải pháp 1 (DB Lock):** Sử dụng `SELECT ... FOR UPDATE` (Pessimistic Lock). An toàn nhưng chậm do tắc nghẽn Database.
- **Giải pháp 2 (Redis Lua - Final):** Chuyển logic trừ kho lên RAM (Redis) và dùng **Lua Script** để đảm bảo tính nguyên tử (Atomicity).
- **Kết quả:** Tăng tốc độ xử lý từ 50 req/s lên **10.000+ req/s**.

### 2. Kiến trúc Microservices & Observability

- **Vấn đề:** Khó quản lý khi hệ thống lớn, REST API chậm chạp. Khó debug khi request đi qua nhiều service.
- **Giải pháp:**
  - Sử dụng **gRPC** (Protobuf) để giao tiếp nội bộ (nhanh gấp 5-10 lần JSON).
  - Triển khai **API Gateway** làm chốt chặn bảo mật (Auth Middleware).
  - Tích hợp **Jaeger & OpenTelemetry** để vẽ biểu đồ Distributed Tracing, giúp phát hiện nút thắt cổ chai (Bottleneck).

### 3. Hệ thống Chat phân tán (Distributed Chat)

- **Vấn đề:** WebSocket chỉ kết nối user với 1 server cụ thể. Khi scale lên 2 server, user ở server khác nhau không chat được.
- **Giải pháp:** Sử dụng **Redis Pub/Sub** làm trung gian chuyển phát tin nhắn giữa các node server.

---

## 🛠 Cài đặt & Chạy thử (Installation)

Dự án được đóng gói hoàn toàn bằng **Docker Compose**. Chỉ cần 1 lệnh để khởi động toàn bộ hệ sinh thái.

### Yêu cầu

- Docker & Docker Compose
- Go 1.22+ (Nếu muốn chạy local)

### Chạy hệ thống

```bash
# 1. Clone repo
git clone [https://github.com/username/go-backend-mastery.git](https://github.com/username/go-backend-mastery.git)
cd go-backend-mastery

# 2. Khởi động Microservices (Gateway, Auth, Jaeger, Redis...)
docker-compose -f docker-compose-microservices.yml up -d

# 3. Khởi động Chat System
docker-compose -f docker-compose-chat.yml up -d
```
