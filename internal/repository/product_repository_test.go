package repository_test

import (
	"context"
	"ecommerce/internal/domain"
	"ecommerce/internal/repository"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
)

// Embed BaseRepositoryTest
type ProductRepositoryTestSuite struct {
	BaseRepositoryTestSuite
	repo domain.ProductRepository
}

// Setup Test (Running Container)
func (ts *ProductRepositoryTestSuite) SetupSuite() {
	ts.BaseRepositoryTestSuite.SetupSuite()

	ts.repo = repository.NewProductRepository(ts.DbPool)
}

func (ts *ProductRepositoryTestSuite) TestCreateProduct() {
	ctx := context.Background()

	tests := []struct {
		name        string
		input       domain.Product
		expectError bool
		expectedErr error
	}{
		{
			name: "Skenario 1: Sukses membuat produk baru",
			input: domain.Product{
				SKU:   "SKU-ARSITEK-002",
				Name:  "Meja Gambar Arsitek",
				Price: decimal.NewFromInt(2500000),
				Stock: 15,
			},
			expectError: false,
			expectedErr: nil,
		},
		{
			name: "Skenario 2: Gagal duplikasi SKU",
			input: domain.Product{
				SKU:   "SKU-ARSITEK-001",
				Name:  "Meja Gambar Tiruan",
				Price: decimal.NewFromInt(1200000),
				Stock: 5,
			},
			expectError: true,
			expectedErr: domain.ErrSKUDuplicate,
		},
	}

	ts.repo.CreateProduct(ctx, domain.Product{SKU: "SKU-ARSITEK-001", Name: "Asli", Price: decimal.NewFromInt(100), Stock: 1})

	for _, tt := range tests {
		ts.Run(tt.name, func() {
			_, err := ts.repo.CreateProduct(ctx, tt.input)
			if tt.expectError {
				ts.ErrorIs(err, tt.expectedErr)
			} else {
				ts.NoError(err)
			}
		})
	}
}

func TestProductRepositorySuite(t *testing.T) {
	suite.Run(t, new(ProductRepositoryTestSuite))
}
