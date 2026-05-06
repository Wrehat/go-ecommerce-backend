package usecase

import (
	"context"
	"ecommerce/internal/domain"
)

// Definisikan struktur untuk ProductUseCase
type ProductUsecase struct {
	// Menyimpan referensi ke repository yang akan digunakan untuk operasi data
	repo domain.ProductRepository
}

// Konstruktor untuk membuat instance baru dari ProductUseCase
func NewProductUsecase(repo domain.ProductRepository) *ProductUsecase {
	return &ProductUsecase{
		repo: repo,
	}
}

// Method untuk menambahkan produk baru
func (u *ProductUsecase) CreateProduct(ctx context.Context, product domain.Product) (domain.Product, error) {
	createdProduct, err := u.repo.CreateProduct(ctx, product)
	if err != nil {
		return domain.Product{}, err
	}

	return createdProduct, nil
}

// Method untuk mengambil semua produk
func (u *ProductUsecase) GetAllProducts(ctx context.Context) ([]domain.Product, error) {
	return u.repo.GetAllProducts(ctx)
}

func (u *ProductUsecase) GetProductByID(ctx context.Context, id int) (domain.Product, error) {
	return u.repo.GetProductByID(ctx, id)
}

// Method update product
func (u *ProductUsecase) UpdateProduct(ctx context.Context, id int, payload domain.Product) (domain.Product, error) {
	return u.repo.UpdateProduct(ctx, id, payload)
}

// Method untuk menghapus produk berdasarkan ID
func (u *ProductUsecase) DeleteProduct(ctx context.Context, id int) error {
	return u.repo.DeleteProduct(ctx, id)
}
