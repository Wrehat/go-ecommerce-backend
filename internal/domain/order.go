package domain

import (
	"context"
	"errors"
	"time"

	"github.com/shopspring/decimal"
)

type Order struct {
	ID          int
	SecureID    string
	UserID      int
	TotalAmount decimal.Decimal
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type OrderItem struct {
	ID        int
	SecureID  string
	OrderID   int
	ProductID int
	Quantity  int
	Price     decimal.Decimal
	Subtotal  decimal.Decimal
}

type CheckOutItemParam struct {
	ProductID int
	Quantity  int
}

type CheckOutParam struct {
	UserID int
	Items  []CheckOutItemParam
}

type OrderRepository interface {
	CreateOrder(ctx context.Context, order Order, items []OrderItem) error
}

type OrderUsecase interface {
	CheckOut(ctx context.Context, param CheckOutParam) (Order, error)
}

var (
	ErrInsufficientStock = errors.New("stock product tidak mencukupi")
)
