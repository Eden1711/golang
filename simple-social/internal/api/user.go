package api

import (
	"simple-social/internal/service"

	"github.com/gin-gonic/gin"
)

func (s *Server) registerUser(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required,alphanum"`
		Password string `json:"password" binding:"required,min=6"`
		Fullname string `json:"fullname" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, err := s.userService.Register(c, service.RegisterReq{
		Username: req.Username,
		Password: req.Password,
		Fullname: req.Fullname,
		Email:    req.Email,
	})

	if err != nil {
		// Thực tế nên check lỗi duplicate key (trùng username) để báo rõ hơn
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, user)

}

func (s *Server) loginUser(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	res, err := s.userService.Login(c, service.LoginReq{
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()}) // 401 Unauthorized
		return
	}

	c.JSON(200, res)
}

func (s *Server) getProfile(c *gin.Context) {
	// 1. Lấy username từ Middleware (Cái mà ta đã c.Set("username_from_token", ...) lúc nãy)
	// MustGet: Nếu không có nó sẽ Panic (nhưng vì đã qua Auth Middleware nên chắc chắn 100% là có)
	username := c.MustGet("username_from_token").(string)

	// 2. Gọi Service lấy thông tin chi tiết
	user, err := s.userService.GetUser(c, username)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// 3. Trả về
	c.JSON(200, user)
}

func (s *Server) followUser(c *gin.Context) {
	var req struct {
		FollowingID int64 `json:"following_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 2. Lấy username từ Token (Middleware)
	username := c.MustGet("username_from_token").(string)

	// 3. Nhưng DB cần ID, không phải Username.
	// Ta phải tìm xem người đang đăng nhập này ID là bao nhiêu.
	// (Thực tế: Nên lưu UserID vào Token luôn để đỡ phải query bước này, nhưng tạm thời làm thế này cho dễ hiểu)
	currentUser, err := s.userService.GetUser(c, username)
	if err != nil {
		c.JSON(500, gin.H{"error": "Lỗi hệ thống khi tìm user"})
		return
	}

	// 4. Gọi Service Follow
	err = s.userService.FollowUser(c, service.FollowUserReq{
		FollowerID:  currentUser.ID,  // Tôi (người đang đăng nhập)
		FollowingID: req.FollowingID, // Người tôi muốn follow
	})

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Follow thành công!"})
}
