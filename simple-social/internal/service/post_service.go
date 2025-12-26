package service

import (
	"context"
	"errors"
	"simple-social/internal/db"
)

type CreatePostReq struct {
	Username string
	Content  string
}

type GetNewsFeedReq struct {
	UserID   int64 // ID của người xem (Me)
	PageID   int32 // Trang số mấy (1, 2, 3...)
	PageSize int32 // Số bài trên 1 trang (10, 20...)
}

type PostService interface {
	CreatePost(ctx context.Context, req CreatePostReq) (db.Post, error)
	GetNewsFeed(ctx context.Context, req GetNewsFeedReq) ([]db.GetNewsFeedRow, error)
}

type postService struct {
	store *db.Queries
}

func NewPostService(store *db.Queries) PostService {
	return &postService{store: store}
}

func (s *postService) CreatePost(ctx context.Context, req CreatePostReq) (db.Post, error) {
	// 1. Phải tìm xem ông User này ID là bao nhiêu?
	user, err := s.store.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return db.Post{}, errors.New("user not found")
	}

	// 2. Tạo bài viết
	arg := db.CreatePostParams{
		UserID:  user.ID,
		Content: req.Content,
	}

	return s.store.CreatePost(ctx, arg)
}

func (s *postService) GetNewsFeed(ctx context.Context, req GetNewsFeedReq) ([]db.GetNewsFeedRow, error) {
	arg := db.GetNewsFeedParams{
		FollowerID: req.UserID,
		Limit:      req.PageSize,
		Offset:     (req.PageID - 1) * req.PageSize, // Công thức tính Offset kinh điển
	}

	return s.store.GetNewsFeed(ctx, arg)
}
