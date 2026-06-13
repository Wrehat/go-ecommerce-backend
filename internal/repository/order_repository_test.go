package repository_test

import (
	"context"
	"ecommerce/internal/domain"
	"ecommerce/internal/repository"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/stretchr/testify/suite"
)

// Definisikan Struct Suite
type OrderRepositoryTestSuite struct {
	BaseRepositoryTestSuite
	userRepo  domain.UserRepository
	prodRepo  domain.ProductRepository
	orderRepo domain.OrderRepository
}

// Jalankan Setup Suite
func (ts *OrderRepositoryTestSuite) SetupSuite() {
	ts.BaseRepositoryTestSuite.SetupSuite()

	ts.userRepo = repository.NewUserRepository(ts.DbPool)
	ts.prodRepo = repository.NewProductRepository(ts.DbPool)
	ts.orderRepo = repository.NewOrderRepository(ts.DbPool)
}

// Panggil Test Suite
func TestOrderRepositorySuite(t *testing.T) {
	suite.Run(t, new(OrderRepositoryTestSuite))
}

// Siapkan Data Master
func (ts *OrderRepositoryTestSuite) SetupMasterData(ctx context.Context) (domain.User, domain.Product) {
	user := domain.User{
		Name:         "Radya",
		Email:        "radya@test.com",
		PasswordHash: "dummy_hashed_password_123",
		Role:         "user",
	}

	createdUser, err := ts.userRepo.CreateUser(ctx, user)
	ts.Require().NoError(err, "Gagal membuat data master user")

	product := domain.Product{
		SKU:   "LAPTOP-001",
		Name:  "Laptop Pro",
		Price: decimal.NewFromInt(10000000),
		Stock: 10,
	}

	createdProduct, err := ts.prodRepo.CreateProduct(ctx, product)
	ts.Require().NoError(err, "Gagal membuat Master Data Product")

	return createdUser, createdProduct
}

// Test 1: Order berhasil, stock berkurang
func (ts *OrderRepositoryTestSuite) TestCreateOrder_Success_StockDecrement() {
	ctx := context.Background()
	user, prod := ts.SetupMasterData(ctx)

	orderParam := domain.Order{
		SecureID:    uuid.New().String(),
		UserID:      user.ID,
		TotalAmount: decimal.NewFromInt(20000000),
		Status:      "PENDING",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	items := []domain.OrderItem{
		{
			SecureID:  uuid.New().String(),
			ProductID: prod.ID,
			Quantity:  2,
			Price:     decimal.NewFromInt(10000000),
		},
	}

	err := ts.orderRepo.CreateOrder(ctx, orderParam, items)
	ts.Require().NoError(err)

	updatedProduct, err := ts.prodRepo.GetProductByID(ctx, prod.ID)
	ts.Require().NoError(err)
	ts.Equal(8, updatedProduct.Stock)
}

func (ts *OrderRepositoryTestSuite) TestCreateOrder_InvalidProduct_OrderNotPersisted() {
	ctx := context.Background()
	user, prod := ts.SetupMasterData(ctx)

	orderParam := domain.Order{
		SecureID:    uuid.New().String(),
		UserID:      user.ID,
		TotalAmount: decimal.NewFromInt(50000000),
		Status:      "PENDING",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	items := []domain.OrderItem{
		{
			SecureID:  uuid.New().String(),
			ProductID: 9999,
			Quantity:  5,
			Price:     decimal.NewFromInt(10000000),
		},
	}

	err := ts.orderRepo.CreateOrder(ctx, orderParam, items)
	ts.Error(err)

	var orderCount int
	errQuery := ts.DbPool.QueryRow(ctx, "SELECT COUNT(*) FROM orders WHERE total_amount = 50000000").Scan(&orderCount)

	ts.Require().NoError(errQuery)
	ts.Equal(0, orderCount, "Data bocor! Transaksi tidak di-rollback dengan benar")

	updatedProduct, err := ts.prodRepo.GetProductByID(ctx, prod.ID)
	ts.Require().NoError(err)
	ts.Equal(10, updatedProduct.Stock)

}
