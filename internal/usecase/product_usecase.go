package usecase

import (
	"context"
	"ecommerce/internal/domain"
	"ecommerce/pkg/logger"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// Definisikan struktur untuk ProductUseCase
type productUsecase struct {
	// Menyimpan referensi ke repository yang akan digunakan untuk operasi data
	repo        domain.ProductRepository
	redisClient *redis.Client
}

// Konstruktor untuk membuat instance baru dari ProductUseCase
func NewProductUsecase(rp domain.ProductRepository, rds *redis.Client) domain.ProductUsecase {
	return &productUsecase{
		repo:        rp,
		redisClient: rds,
	}
}

// Method untuk menambahkan produk baru
func (u *productUsecase) CreateProduct(ctx context.Context, product domain.Product) (domain.Product, error) {
	createdProduct, err := u.repo.CreateProduct(ctx, product)
	if err != nil {
		return domain.Product{}, err
	}

	u.redisClient.Del(ctx, "product:all")

	return createdProduct, nil
}

// Method untuk mengambil semua produk
func (u *productUsecase) GetAllProducts(ctx context.Context) ([]domain.Product, error) {
	// Todo : buat key untuk ambil data redis
	cacheKey := "product:all"

	// Todo : Usecase cek ke redis
	cachedData, err := u.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		fmt.Println("🚀 CACHE HIT: Ambil dari Redis!")
		var product []domain.Product
		if err := json.Unmarshal([]byte(cachedData), &product); err != nil {
			logger.Log.Error("Error Unmarshal data cache", zap.Error(err))
		}
		return product, nil
	}

	// Todo : Usecase ambil dari database
	fmt.Println("🐢 CACHE MISS: Ambil dari Postgres!")
	products, err := u.repo.GetAllProducts(ctx)
	if err != nil {
		return nil, err
	}

	// Todo : Simpan data ke redis
	productsJSON, _ := json.Marshal(products)
	if err := u.redisClient.Set(ctx, cacheKey, productsJSON, 5*time.Minute).Err(); err != nil {
		logger.Log.Error("Error Saat Simpan Ke Redis", zap.Error(err))
	}

	return products, nil
}

func (u *productUsecase) GetProductByID(ctx context.Context, id int) (domain.Product, error) {
	return u.repo.GetProductByID(ctx, id)
}

// Method update product
func (u *productUsecase) UpdateProduct(ctx context.Context, id int, payload domain.Product) (domain.Product, error) {
	u.redisClient.Del(ctx, "product:all")
	return u.repo.UpdateProduct(ctx, id, payload)
}

// Method untuk menghapus produk berdasarkan ID
func (u *productUsecase) DeleteProduct(ctx context.Context, id int) error {
	u.redisClient.Del(ctx, "product:all")
	return u.repo.DeleteProduct(ctx, id)
}
