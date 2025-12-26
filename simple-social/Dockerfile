# Giai đoạn 1: Build (Dùng ảnh chứa Go để compile code)
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
# Build ra file thực thi tên là 'main'
RUN go build -o main cmd/main.go

# Giai đoạn 2: Run (Dùng ảnh Alpine siêu nhẹ để chạy)
FROM alpine:3.19
WORKDIR /app
# Copy file thực thi từ giai đoạn 1 sang giai đoạn 2
COPY --from=builder /app/main .
# COPY app.env . 
# (Lưu ý: Nếu bạn dùng file config app.env thì nhớ copy, nếu không thì bỏ dòng trên)

# Mở cổng 8080
EXPOSE 8080
# Lệnh chạy khi bật container
CMD ["/app/main"]