package service

import (
	"context"
	"errors"
	"simple-bank/internal/db"
)

type CreateAccountReq struct {
	Owner    string
	Currency string
}

type WithdrawReq struct {
	AccountID int64
	Amount    int64
}

type TransferReq struct {
	FromAccountID int64
	ToAccountID   int64
	Amount        int64
}

type AccountService interface {
	CraeteAccount(ctx context.Context, req CreateAccountReq) (db.Account, error)
	WithdrawMoney(ctx context.Context, req WithdrawReq) (db.Account, error)
	TransferMoney(ctx context.Context, req TransferReq) (db.Transfer, error)
}

type accountService struct {
	store *db.Queries
}

func NewAccountService(store *db.Queries) AccountService {
	return &accountService{store: store}
}

// Logic 1: Tạo tài khoản (Tặng sẵn 0 đồng)
func (s *accountService) CraeteAccount(ctx context.Context, req CreateAccountReq) (db.Account, error) {
	if req.Owner == "" {
		return db.Account{}, errors.New("tên chủ tài khoản không được trống")
	}

	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0, // Mặc định 0 đồng
	}
	return s.store.CreateAccount(ctx, arg)
}

// Logic 2: Rút tiền (Logic nghiệp vụ nằm ở đây!)
func (s *accountService) WithdrawMoney(ctx context.Context, req WithdrawReq) (db.Account, error) {
	// B1: Lấy thông tin tài khoản hiện tại
	account, err := s.store.GetAccount(ctx, req.AccountID)
	if err != nil {
		return db.Account{}, errors.New("tài khoản không tồn tại")
	}

	// B2: Kiểm tra số dư (Business Rule)
	if account.Balance < req.Amount {
		return db.Account{}, errors.New("số dư không đủ để thực hiện giao dịch")
	}

	// B3: Thực hiện trừ tiền
	// Truyền số âm vào để trừ
	arg := db.UpdateAccountBalanceParams{
		ID:     req.AccountID,
		Amount: -req.Amount,
	}

	return s.store.UpdateAccountBalance(ctx, arg)
}

func (s *accountService) TransferMoney(ctx context.Context, req TransferReq) (db.Transfer, error) {
	// --- PHẦN 1: VALIDATION DỮ LIỆU ---

	// 1. Kiểm tra chuyển cho chính mình?
	if req.FromAccountID == req.ToAccountID {
		return db.Transfer{}, errors.New("không thể tự chuyển tiền cho chính mình")
	}

	// 2. Lấy thông tin người gửi
	fromAcc, err := s.store.GetAccount(ctx, req.FromAccountID)
	if err != nil {
		return db.Transfer{}, errors.New("tài khoản nguồn không tồn tại")
	}

	// 3. Lấy thông tin người nhận
	toAcc, err := s.store.GetAccount(ctx, req.ToAccountID)
	if err != nil {
		return db.Transfer{}, errors.New("tài khoản đích không tồn tại")
	}

	// 4. Kiểm tra lệch tiền tệ (Logic nghiệp vụ)
	if fromAcc.Currency != toAcc.Currency {
		return db.Transfer{}, errors.New("tiền tệ không khớp (ví dụ: không thể chuyển USD sang VND)")
	}

	// 5. Kiểm tra số dư
	if fromAcc.Balance < req.Amount {
		return db.Transfer{}, errors.New("số dư không đủ")
	}

	// --- PHẦN 2: XỬ LÝ DỮ LIỆU (EXECUTION) ---
	// (Lưu ý: Trong thực tế đoạn này cần Database Transaction - TX,
	// nhưng để học logic ta tạm làm từng bước)

	// Bước A: Trừ tiền người gửi
	_, err = s.store.UpdateAccountBalance(ctx, db.UpdateAccountBalanceParams{
		ID:     req.FromAccountID,
		Amount: -req.Amount, // Số âm
	})
	if err != nil {
		return db.Transfer{}, err
	}

	// Bước B: Cộng tiền người nhận
	_, err = s.store.UpdateAccountBalance(ctx, db.UpdateAccountBalanceParams{
		ID:     req.ToAccountID,
		Amount: req.Amount, // Số dương
	})
	if err != nil {
		// Nguy hiểm: Nếu bước A xong mà bước B lỗi -> Mất tiền oan!
		// (Bài sau ta sẽ học cách fix chỗ này bằng Transaction)
		return db.Transfer{}, err
	}

	// Bước C: Lưu lịch sử giao dịch
	arg := db.CreateTransferParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}
	transfer, err := s.store.CreateTransfer(ctx, arg)
	if err != nil {
		return db.Transfer{}, err
	}

	return transfer, nil
}
