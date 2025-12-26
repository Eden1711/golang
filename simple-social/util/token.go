package util

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(username string, duration time.Duration, secretKey string) (string, error) {
	// 1.Tạo payload (Claims) - Chứa thông tin user
	claims := jwt.MapClaims{"username": username, "exp": time.Now().Add(duration).Unix(), "iat": time.Now().Unix()}

	// 2. Ký tên (Sign) bằng thuật toán HS256 và Secret Key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// VerifyToken kiểm tra token có hợp lệ không? Nếu đúng trả về payload (data bên trong)
func VerifyToken(tokenString string, secretKey string) (jwt.MapClaims, error) {
	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Kiểm tra thuật toán ký có đúng là HMAC không (quan trọng!)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	// Lấy dữ liệu bên trong (Claims)
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
