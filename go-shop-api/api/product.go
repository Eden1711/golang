package api

import (
	"fmt"
	"net/http"
	"strconv"

	"go-shop-api/db" // Import struct params từ SQLC

	"github.com/gin-gonic/gin"
)

// Chú ý: (server *Server) đứng trước tên hàm
// Đây gọi là Method Receiver. Giúp hàm này truy cập được server.store
func (server *Server) createProduct(ctx *gin.Context) {
	var req struct {
		Name  string `json:"name" binding:"required"`
		Sku   string `json:"sku" binding:"required"`
		Price int64  `json:"price" binding:"required,gt=0"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	arg := db.CreateProductParams{
		Name:  req.Name,
		Sku:   req.Sku,
		Price: req.Price,
	}

	// Gọi store thông qua biến server (Dependency Injection thành công!)
	product, err := server.store.CreateProduct(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, product)
}

func (server *Server) listProducts(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	arg := db.ListProductsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	products, err := server.store.ListProducts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, products)
}

func (server *Server) getProduct(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	product, err := server.store.GetProduct(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	ctx.JSON(http.StatusOK, product)
}

func (server *Server) updateProduct(ctx *gin.Context) {
	var req struct {
		Name  string `json:"name" binding:"required"`
		Price int64  `json:"price" binding:"required,gt=0"`
	}

	// 1. Lấy ID từ URL
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// 2. Validate dữ liệu JSON gửi lên
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. Gọi SQLC để Update
	// Lưu ý: arg phải khớp với struct UpdateProductParams trong file db/models.go
	arg := db.UpdateProductParams{
		ID:    id,
		Name:  req.Name,
		Price: req.Price,
	}

	err = server.store.UpdateProduct(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update thành công (Trả về thông báo hoặc object rỗng đều được)
	ctx.JSON(http.StatusOK, gin.H{"status": "updated successfully"})
}

func (server *Server) deleteProduct(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	fmt.Println(id, err)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	_, err = server.store.GetProduct(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Not Found"})
		return
	}

	err = server.store.DeleteProduct(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 3. Trả về thành công
	ctx.JSON(http.StatusOK, gin.H{"status": "deleted successfully"})

}
