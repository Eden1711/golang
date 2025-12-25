package api

import (
	"net/http"
	"simple-bank/internal/service"

	"github.com/gin-gonic/gin"
)

type Server struct {
	service service.AccountService
	Router  *gin.Engine
}

func NewServer(svc service.AccountService) *Server {
	server := &Server{service: svc}
	server.setupRouter()
	return server
}

func (s *Server) setupRouter() {
	r := gin.Default()
	r.POST("/accounts", s.createAccount)
	r.POST("/withdraw", s.withdraw)      // API rút tiền
	r.POST("/transfer", s.transferMoney) // API rút tiền
	s.Router = r
}

func (s *Server) createAccount(c *gin.Context) {
	var req struct {
		Owner    string `json:"owner" binding:"required"`
		Currency string `json:"currency" binding:"required,oneof=USD EUR VND"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	svcReq := service.CreateAccountReq{
		Owner:    req.Owner,
		Currency: req.Currency,
	}

	acc, err := s.service.CraeteAccount(c, svcReq)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, acc)
}

func (s *Server) withdraw(c *gin.Context) {
	var req struct {
		AccountID int64 `json:"account_id" binding:"required"`
		Amount    int64 `json:"amount" binding:"required,gt=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	svcReq := service.WithdrawReq{
		AccountID: req.AccountID,
		Amount:    req.Amount,
	}

	// Gọi Service để rút
	newAcc, err := s.service.WithdrawMoney(c, svcReq)
	if err != nil {
		// Ở đây có thể check lỗi để trả về 400 (Bad Request) nếu không đủ tiền
		// Nhưng để đơn giản ta trả 500
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"message":     "Rút tiền thành công",
		"new_balance": newAcc.Balance,
	})
}

func (s *Server) transferMoney(c *gin.Context) {
	var req struct {
		FromAccountID int64 `json:"from_account_id" binding:"required"`
		ToAccountID   int64 `json:"to_account_id" binding:"required"`
		Amount        int64 `json:"amount" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	svcReq := service.TransferReq{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := s.service.TransferMoney(c, svcReq)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Chuyển tiền thành công",
		"data":    result,
	})
}
