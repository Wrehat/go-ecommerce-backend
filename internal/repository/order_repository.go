package repository

import (
	"context"
	"ecommerce/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type orderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) domain.OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) CreateOrder(ctx context.Context, order domain.Order, items []domain.OrderItem) error {
	// Todo: mulai transaksi (Begin)
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	// Todo: pasang jaring pengaman
	defer tx.Rollback(ctx)

	// Todo: Simpan data order, returning id untuk order item
	var orderID int
	orderQuery := `INSERT INTO orders (secure_id, user_id, total_amount, status, created_at, updated_at) VALUES ($1,$2,$3,$4, $5, $6) RETURNING id`

	err = tx.QueryRow(ctx, orderQuery, order.SecureID, order.UserID, order.TotalAmount, order.Status, order.CreatedAt, order.UpdatedAt).Scan(&orderID)
	if err != nil {
		return err
	}

	// Todo: proses item order
	for _, item := range items {
		var currentStock int

		// Todo: lock read db
		checkStockQuery := `SELECT stock FROM products WHERE id=$1 AND deleted_at IS NULL FOR UPDATE`

		err := tx.QueryRow(ctx, checkStockQuery, item.ProductID).Scan(&currentStock)
		if err != nil {
			return domain.ErrProductNotFound
		}

		// Todo: Cek currentStock < item.Quantity
		if currentStock < item.Quantity {
			return domain.ErrInsufficientStock
		}

		// Todo: Kurangi stock
		updateStockQuery := `UPDATE products SET stock = stock - $1 WHERE id=$2`

		_, err = tx.Exec(ctx, updateStockQuery, item.Quantity, item.ProductID)
		if err != nil {
			return err
		}

		// Todo: Simpan data order item
		itemQuery := `INSERT INTO order_items (secure_id, order_id, product_id, quantity, price) VALUES ($1,$2,$3,$4,$5)`

		_, err = tx.Exec(ctx, itemQuery, item.SecureID, orderID, item.ProductID, item.Quantity, item.Price)
		if err != nil {
			return err
		}
	}

	// Sahkan commit
	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
