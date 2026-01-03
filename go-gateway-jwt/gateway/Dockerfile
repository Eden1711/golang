# Giai đoạn 1: Builder
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Build file gateway
RUN go build -o gatewayApp ./gateway/main.go

# Giai đoạn 2: Runner
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/gatewayApp .

EXPOSE 8080
CMD ["./gatewayApp"]