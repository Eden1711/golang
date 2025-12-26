package api

import (
	"simple-social/internal/service"

	"github.com/gin-gonic/gin"
)

type Server struct {
	userService service.UserService
	postService service.PostService
	Router      *gin.Engine
}

// NewServer khởi tạo Server và Setup Router
func NewServer(userSvc service.UserService, postSvc service.PostService) *Server {
	server := &Server{
		userService: userSvc, // Gán Service vào struct
		postService: postSvc,
	}

	server.setupRouter()
	return server
}

func (s *Server) setupRouter() {
	r := gin.Default()
	v1 := r.Group("/api/v1")
	{
		v1.POST("/register", s.registerUser)
		v1.POST("/login", s.loginUser)
		// 👇 Nhóm API cần bảo vệ
		authRoutes := v1.Group("/").Use(authMiddleware())
		{
			authRoutes.POST("/posts", s.createPost)
			// 👇 Thêm dòng này: GET /users/me
			authRoutes.GET("/users/me", s.getProfile)
			authRoutes.POST("/follow", s.followUser)
			authRoutes.GET("/posts/feed", s.getNewsFeed)
		}
	}

	s.Router = r
}
