package repository

import (
	"context"
	"ecommerce/internal/domain"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type productRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) domain.ProductRepository {
	return &productRepository{db: db}
}

// Method untuk mengambil semua produk dari repository
func (r *productRepository) GetAllProducts(ctx context.Context) ([]domain.Product, error) {
	// 1. Buat slice untuk menyimpan produk
	products := make([]domain.Product, 0)
	// 2. Buat query untuk mengambil semua produk dari database
	query := "SELECT id, sku, name, price, stock, created_at, updated_at FROM products WHERE deleted_at IS NULL"
	// 3. Eksekusi query
	rows, err := r.db.Query(ctx, query)
	// 4. Cek error saat eksekusi query
	if rows == nil || err != nil {
		return nil, fmt.Errorf("Gagal mengambil produk: %v", err)
	}
	// 5. Pastikan rows ditutup setelah selesai digunakan
	defer rows.Close()
	// 6. Iterasi hasil query dan masukkan ke dalam slice products
	for rows.Next() {
		// 6a. Buat variabel untuk menyimpan data produk sementara
		var product domain.Product
		// 6b. Scan data dari row ke var product
		err = rows.Scan(&product.ID, &product.SKU, &product.Name, &product.Price, &product.Stock, &product.CreatedAt, &product.UpdatedAt)
		// 6c. Cek error saat scanning data
		if err != nil {
			return nil, fmt.Errorf("Gagal memproses data produk: %v", err)
		}
		// 6d. Tambahkan produk ke dalam slice products
		products = append(products, product)
	}
	// 7. Cek error setelah iterasi selesai
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("Error iterasi hasil query: %v", err)
	}
	// 8. Kembalikan slice products dan error (jika ada)
	return products, nil
}

// Method untuk membuat produk dari repository
func (r *productRepository) CreateProduct(ctx context.Context, product domain.Product) (domain.Product, error) {
	// Buat Query
	query := `INSERT INTO products (sku,name,price,stock) VALUES ($1,$2,$3,$4) RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(ctx, query, product.SKU, product.Name, product.Price, product.Stock).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.Product{}, domain.ErrSKUDuplicate
		}

		return domain.Product{}, err
	}

	return product, nil

}

func (r *productRepository) GetProductByID(ctx context.Context, id int) (domain.Product, error) {
	// Todo : Buat query dan wadah
	query := `SELECT id, sku, name, price, stock, created_at, updated_at from products WHERE id= $1 AND deleted_at IS NULL`
	var p domain.Product

	// Todo : panggil dan error handling
	err := r.db.QueryRow(ctx, query, id).Scan(&p.ID, &p.SKU, &p.Name, &p.Price, &p.Stock, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Product{}, domain.ErrProductNotFound
		}
		return domain.Product{}, err
	}

	// Todo : return response
	return p, nil
}

func (r *productRepository) UpdateProduct(ctx context.Context, id int, payload domain.Product) (domain.Product, error) {

	// Todo: buat query update
	query := `UPDATE products SET SKU=$1, name=$2, price=$3, stock=$4, updated_at= CURRENT_TIMESTAMP WHERE id=$5 AND deleted_at IS NULL RETURNING id, SKU, name, price, stock, updated_at`

	row := r.db.QueryRow(ctx, query, payload.SKU, payload.Name, payload.Price, payload.Stock, id)

	var p domain.Product
	err := row.Scan(&p.ID, &p.SKU, &p.Name, &p.Price, &p.Stock, &p.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Product{}, domain.ErrProductNotFound
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return domain.Product{}, domain.ErrSKUDuplicate
			}
		}
		return domain.Product{}, err
	}
	return p, nil

}

func (r *productRepository) DeleteProduct(ctx context.Context, id int) error {
	// Todo: Buat query delete
	query := `UPDATE products SET deleted_at=CURRENT_TIMESTAMP WHERE id=$1 AND deleted_at IS NULL`

	// Todo: Jalankan query dan handle error
	commandTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return domain.ErrProductNotFound
	}

	// Todo: return success response
	return nil
}
