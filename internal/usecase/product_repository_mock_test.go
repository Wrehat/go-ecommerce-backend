package usecase_test

// import (
// 	"context"

// 	"ecommerce/internal/domain"

// 	"github.com/stretchr/testify/mock"
// )

// type ProductRepositoryMock struct {
// 	mock.Mock
// }

// func (m *ProductRepositoryMock) GetProductByID(ctx context.Context, id int) (domain.Product, error) {
// 	args := m.Called(ctx, id)
// 	return args.Get(0).(domain.Product), args.Error(1)
// }

// func (m *ProductRepositoryMock) CreateProduct(ctx context.Context, p domain.Product) (domain.Product, error) {
// 	args := m.Called(ctx, p)

// 	return args.Get(0).(domain.Product), args.Error(0)
// }

// func (m *ProductRepositoryMock) DeleteProduct(ctx context.Context, id int) error {
// 	args := m.Called(ctx, id)

// 	return args.Error(0)
// }

// func (m *ProductRepositoryMock) GetAllProducts(ctx context.Context) ([]domain.Product, error) {
// 	args := m.Called(ctx)

// 	return args.Get(0).([]domain.Product), args.Error(0)
// }

// func (m *ProductRepositoryMock) UpdateProduct(ctx context.Context, id int, p domain.Product) (domain.Product, error) {
// 	args := m.Called(ctx, id, p)

// 	return args.Get(0).(domain.Product), args.Error(0)

// }
