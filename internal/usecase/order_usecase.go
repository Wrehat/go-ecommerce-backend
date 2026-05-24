package usecase

import (
	"context"
	"ecommerce/internal/domain"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type orderUsecase struct {
	orderRepo   domain.OrderRepository
	productRepo domain.ProductRepository
}

func NewOrderUsecase(or domain.OrderRepository, pr domain.ProductRepository) domain.OrderUsecase {
	return &orderUsecase{
		orderRepo:   or,
		productRepo: pr,
	}
}

func (u *orderUsecase) CheckOut(ctx context.Context, param domain.CheckOutParam) (domain.Order, error) {
	order := domain.Order{
		SecureID:    uuid.New().String(),
		UserID:      param.UserID,
		Status:      "PENDING",
		TotalAmount: decimal.NewFromInt(0),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	var orderItems []domain.OrderItem

	for _, itemParam := range param.Items {
		p, err := u.productRepo.GetProductByID(ctx, itemParam.ProductID)
		if err != nil {
			return domain.Order{}, err
		}

		qtyDecimal := decimal.NewFromInt(int64(itemParam.Quantity))
		subtotal := p.Price.Mul(qtyDecimal)

		orderItem := domain.OrderItem{
			SecureID:  uuid.New().String(),
			ProductID: itemParam.ProductID,
			Quantity:  itemParam.Quantity,
			Price:     p.Price,
			Subtotal:  subtotal,
		}
		orderItems = append(orderItems, orderItem)

		order.TotalAmount = order.TotalAmount.Add(subtotal)
	}

	if err := u.orderRepo.CreateOrder(ctx, order, orderItems); err != nil {
		return domain.Order{}, err
	}

	return order, nil
}
