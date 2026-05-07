package domain

import (
	"context"
	"errors"
	"time"

	"github.com/shopspring/decimal"
)

type Product struct {
	ID        int
	SKU       string
	Name      string
	Price     decimal.Decimal
	Stock     int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ProductRepository interface {
	GetAllProducts(context.Context) ([]Product, error)
	CreateProduct(context.Context, Product) (Product, error)
	GetProductByID(context.Context, int) (Product, error)
	UpdateProduct(context.Context, int, Product) (Product, error)
	DeleteProduct(context.Context, int) error
}

var ErrSKUDuplicate = errors.New("sku sudah terdaftar, silakan gunakan sku lain")
var ErrProductNotFound = errors.New("product tidak ditemukan")
