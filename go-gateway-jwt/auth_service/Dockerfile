# Giai đoạn 1: Builder (Dùng ảnh nặng để compile code)
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Copy file dependency trước để tận dụng cache
COPY go.mod go.sum ./
RUN go mod download

# Copy toàn bộ code và Build
COPY . .
# Build ra file nhị phân tên là 'authApp'
RUN go build -o authApp ./auth_service/main.go

# Giai đoạn 2: Runner (Dùng ảnh siêu nhẹ để chạy)
FROM alpine:latest
WORKDIR /root/
# Chỉ copy file nhị phân từ giai đoạn 1 sang
COPY --from=builder /app/authApp .

# Mở port
EXPOSE 50051
# Chạy app
CMD ["./authApp"]