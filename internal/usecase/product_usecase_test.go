package usecase_test

import (
	"context"
	"ecommerce/internal/domain"
	"ecommerce/internal/usecase"
	"ecommerce/internal/usecase/mocks"

	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestProductUsecase_GetProductByID(t *testing.T) {
	// Ambil context
	ctx := context.Background()

	dummyProd := domain.Product{
		ID:    1,
		SKU:   "SKU001",
		Name:  "Buku",
		Price: decimal.NewFromInt(10000),
		Stock: 10,
	}

	tests := []struct {
		name        string
		inputID     int
		mockSetup   func(r *mocks.MockProductRepository)
		expectError bool
		expectedErr error
	}{
		{
			name:    "Skenario 1: Sukses ambil data",
			inputID: 1,
			mockSetup: func(r *mocks.MockProductRepository) {
				r.On("GetProductByID", ctx, 1).Return(dummyProd, nil).Once()
			},
			expectError: false,
			expectedErr: nil,
		},
		{
			name:    "Skenario 2: Gagal ambil data",
			inputID: 99,
			mockSetup: func(r *mocks.MockProductRepository) {
				r.On("GetProductByID", ctx, 99).Return(domain.Product{}, domain.ErrProductNotFound).Once()
			},
			expectError: true,
			expectedErr: domain.ErrProductNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockProductRepository)

			prodUC := usecase.NewProductUsecase(mockRepo, nil)

			tt.mockSetup(mockRepo)

			prod, err := prodUC.GetProductByID(ctx, tt.inputID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.inputID, prod.ID)
				assert.Equal(t, dummyProd.Name, prod.Name)
				assert.Equal(t, dummyProd.Price, prod.Price)
			}

			mockRepo.AssertExpectations(t)
		})
	}

}
