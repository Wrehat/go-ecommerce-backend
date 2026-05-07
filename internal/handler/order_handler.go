package handler

import (
	"ecommerce/internal/domain"
	"ecommerce/pkg/response"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CheckOutItemReq struct {
	ProductID int `json:"product_id" binding:"required"`
	Quantity  int `json:"quantity" binding:"required,min=1"`
}

type CheckOutReq struct {
	UserID int               `json:"user_id" binding:"required"`
	Items  []CheckOutItemReq `json:"items" binding:"required,min=1,dive"`
}

type OrderResponse struct {
	SecureID    string `json:"secure_id"`
	UserID      int    `json:"user_id"`
	TotalAmount string `json:"total_amount"`
	Status      string `json:"status"`
}

type orderHandler struct {
	usecase domain.OrderUsecase
}

func NewOrderHandler(oh domain.OrderUsecase) *orderHandler {
	return &orderHandler{usecase: oh}
}

func (h *orderHandler) CheckOut(c *gin.Context) {
	// Todo: Bind JSON
	var req CheckOutReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, "Format data tidak valid", nil)
		return
	}

	// Todo: Map DTO ke Param Usecase
	paramItems := make([]domain.CheckOutItemParam, 0, len(req.Items))
	for _, item := range req.Items {
		paramItems = append(paramItems, domain.CheckOutItemParam{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	param := domain.CheckOutParam{
		UserID: req.UserID,
		Items:  paramItems,
	}

	// Todo: panggil usecase
	ctx := c.Request.Context()
	order, err := h.usecase.CheckOut(ctx, param)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			response.JSON(c, http.StatusNotFound, "Product tidak ditemukan", nil)
			return
		}
		if errors.Is(err, domain.ErrInsufficientStock) {
			response.JSON(c, http.StatusBadRequest, "Stock product tidak cukup", nil)
			return
		}
		response.JSON(c, http.StatusInternalServerError, "Terjadi kesalahan pada server", nil)
		return
	}

	res := OrderResponse{
		SecureID:    order.SecureID,
		UserID:      req.UserID,
		TotalAmount: order.TotalAmount.String(),
		Status:      order.Status,
	}

	response.JSON(c, http.StatusCreated, "Berhasil membuat pesanan", res)
}
