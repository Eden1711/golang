package service

import (
	"context"
	"errors"
	"simple-social/internal/db"
	"simple-social/util"
	"time"
)

type RegisterReq struct {
	Username string
	Password string
	Fullname string
	Email    string
}

type LoginReq struct {
	Username string
	Password string
}

type LoginResponse struct {
	AccessToken string
	User        db.User
}

type FollowUserReq struct {
	FollowerID  int64 // ID của người đi follow (Lấy từ Token)
	FollowingID int64 // ID của người được follow (Gửi từ Client)
}

type UserService interface {
	Register(ctx context.Context, req RegisterReq) (db.User, error)
	Login(ctx context.Context, req LoginReq) (LoginResponse, error)
	GetUser(ctx context.Context, username string) (db.User, error)
	FollowUser(ctx context.Context, req FollowUserReq) error
}

type userService struct {
	store *db.Queries
}

func NewUserService(store *db.Queries) UserService {
	return &userService{store: store}
}

func (s *userService) FollowUser(ctx context.Context, req FollowUserReq) error {
	if req.FollowerID == req.FollowingID {
		return errors.New("không thể follow chính mình")
	}

	arg := db.CreateFollowParams{
		FollowerID:  req.FollowerID,
		FollowingID: req.FollowingID,
	}

	_, err := s.store.CreateFollow(ctx, arg)

	if err != nil {
		return err
	}

	return nil
}

func (s *userService) Register(ctx context.Context, req RegisterReq) (db.User, error) {
	hashed, err := util.HashPassword(req.Password)

	if err != nil {
		return db.User{}, err
	}

	arg := db.CreateUserParams{
		Username:     req.Username,
		PasswordHash: hashed,
		FullName:     req.Fullname,
		Email:        req.Email}

	user, err := s.store.CreateUser(ctx, arg)

	if err != nil {
		return db.User{}, err
	}

	user.PasswordHash = ""

	return user, nil
}

func (s *userService) Login(ctx context.Context, req LoginReq) (LoginResponse, error) {
	// 1. Tìm user trong DB
	user, err := s.store.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return LoginResponse{}, errors.New("tài khoản không tồn tại")
	}

	// 2. Kiểm tra mật khẩu (So sánh cái nhập vào vs cái trong DB)
	err = util.CheckPassword(req.Password, user.PasswordHash)

	if err != nil {
		return LoginResponse{}, errors.New("sai mật khẩu")
	}

	// 3. Mật khẩu đúng -> Tạo Token (Vé thông hành)
	// Token sống trong 24 giờ
	token, err := util.CreateToken(user.Username, 24*time.Hour)
	if err != nil {
		return LoginResponse{}, err
	}

	// 4. Trả về kết quả
	user.PasswordHash = "" // Giấu hash đi
	return LoginResponse{
		AccessToken: token,
		User:        user,
	}, nil
}

func (s *userService) GetUser(ctx context.Context, username string) (db.User, error) {
	user, err := s.store.GetUserByUsername(ctx, username)
	if err != nil {
		return db.User{}, err
	}

	// Quan trọng: Phải xóa Hash Password trước khi trả về
	user.PasswordHash = ""
	return user, nil
}
