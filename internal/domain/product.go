package domain

import (
	"context"
	"errors"
	"time"

	"github.com/shopspring/decimal"
)

// type Product struct {
// 	ID        int       `json:"id"`
// 	SKU       string    `json:"sku" binding:"required"`
// 	Name      string    `json:"name" binding:"required"`
// 	Price     decimal.Decimal `json:"price" binding:"required,gt=0"`
// 	Stock     int       `json:"stock" binding:"required,gte=0"`
// 	CreatedAt time.Time `json:"created_at"`
// 	UpdatedAt time.Time `json:"updated_at"`
// }
// TIDAK PAKAI STRUCT TAG BINDING KARENA SUDAH MENERAPKAN VALIDASI DI LAYER HANDLER (PRODUCTREQUEST)

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
