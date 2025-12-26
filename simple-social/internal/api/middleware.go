package api

import (
	"net/http"
	"simple-social/util"
	"strings"

	"github.com/gin-gonic/gin"
)

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Lấy header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			return
		}

		// 2. Format chuẩn phải là: "Bearer <token>"
		fields := strings.Fields(authHeader)
		if len(fields) < 2 || strings.ToLower(fields[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			return
		}

		// 3. Lấy token và Verify
		accessToken := fields[1]
		claims, err := util.VerifyToken(accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		// 4. Token ngon! Lấy username gắn vào context để các hàm sau dùng
		username := claims["username"].(string)
		c.Set("username_from_token", username) // <-- Lưu lại

		c.Next() // Cho phép đi tiếp vào Controller
	}
}
