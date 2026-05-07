package usecase

import (
	"context"
	"ecommerce/internal/domain"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type orderUsecase struct {
	orderRepo   domain.OrderRepository
	productRepo domain.ProductRepository
}

func NewOrderUsecase(or domain.OrderRepository, pr domain.ProductRepository) domain.OrderUsecase {
	return &orderUsecase{
		orderRepo:   or,
		productRepo: pr,
	}
}

func (u *orderUsecase) CheckOut(ctx context.Context, param domain.CheckOutParam) (domain.Order, error) {
	// TODO 1: Siapkan Kertas Nota Kosong (domain.Order)
	// - Generate SecureID pakai: uuid.New().String()
	// - Isi UserID dari param
	// - Set Status jadi "PENDING"
	// - Set TotalAmount jadi 0 pakai: decimal.NewFromInt(0)
	order := domain.Order{
		SecureID:    uuid.New().String(),
		UserID:      param.UserID,
		Status:      "PENDING",
		TotalAmount: decimal.NewFromInt(0),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Siapkan array kosong untuk daftar barang
	var orderItems []domain.OrderItem

	// TODO 2: Lakukan perulangan (loop) untuk setiap barang yang ada di param.Items
	for _, itemParam := range param.Items {
		// TODO 3: Cek harga asli ke database produk menggunakan u.productRepo.GetProductByID
		// (Jika error, return domain.Order{} kosong dan error-nya)
		p, err := u.productRepo.GetProductByID(ctx, itemParam.ProductID)
		if err != nil {
			return domain.Order{}, err
		}

		// TODO 4: Hitung subtotal barang ini (Quantity dikali Harga Asli)
		// Petunjuk konversi: qtyDecimal := decimal.NewFromInt(int64(itemParam.Quantity))
		// subtotal := product.Price.Mul(qtyDecimal)
		qtyDecimal := decimal.NewFromInt(int64(itemParam.Quantity))
		subtotal := p.Price.Mul(qtyDecimal)

		// TODO 5: Buat objek domain.OrderItem dan masukkan ke array orderItems menggunakan append()
		// - Generate SecureID baru untuk item ini
		// - Isi ProductID, Quantity
		// - Isi Price dengan harga asli dari database (Snapshot!)
		// - Isi Subtotal dengan hasil hitungan di TODO 4
		orderItem := domain.OrderItem{
			SecureID:  uuid.New().String(),
			ProductID: itemParam.ProductID,
			Quantity:  itemParam.Quantity,
			Price:     p.Price,
			Subtotal:  subtotal,
		}
		orderItems = append(orderItems, orderItem)

		// TODO 6: Tambahkan subtotal barang ini ke TotalAmount di nota utama (order.TotalAmount)
		// Petunjuk: order.TotalAmount = order.TotalAmount.Add(subtotal)
		order.TotalAmount = order.TotalAmount.Add(subtotal)
	}

	// TODO 7: Serahkan nota dan daftar barang ke Petugas Gudang (u.orderRepo.CreateOrder)
	// (Jika error, return domain.Order{} kosong dan error-nya)
	if err := u.orderRepo.CreateOrder(ctx, order, orderItems); err != nil {
		return domain.Order{}, err
	}

	return order, nil
}
