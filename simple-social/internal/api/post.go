package api

import (
	"simple-social/internal/service"

	"github.com/gin-gonic/gin"
)

func (s *Server) createPost(c *gin.Context) {
	var req struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 🔥 LẤY USERNAME TỪ MIDDLEWARE (Cái mà ta đã c.Set lúc nãy)
	username := c.MustGet("username_from_token").(string)

	post, err := s.postService.CreatePost(c, service.CreatePostReq{
		Username: username,
		Content:  req.Content,
	})

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, post)
}

func (s *Server) getNewsFeed(c *gin.Context) {
	// 1. Lấy tham số phân trang từ URL
	var req struct {
		PageID   int32 `form:"page_id" binding:"required,min=1"`
		PageSize int32 `form:"page_size" binding:"required,min=5,max=50"`
	}

	// ShouldBindQuery: Lấy từ ?page_id=... chứ không phải JSON body
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 2. Xác định "Tôi là ai?"
	username := c.MustGet("username_from_token").(string)

	// Lại phải tìm ID của mình (Hơi cực, bài sau tối ưu chỗ này nhé)
	currentUser, err := s.userService.GetUser(c, username)
	if err != nil {
		c.JSON(500, gin.H{"error": "Không tìm thấy user"})
		return
	}

	// 3. Gọi Service lấy Feed
	posts, err := s.postService.GetNewsFeed(c, service.GetNewsFeedReq{
		UserID:   currentUser.ID,
		PageID:   req.PageID,
		PageSize: req.PageSize,
	})

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// 4. Trả về danh sách bài viết
	c.JSON(200, posts)
}
