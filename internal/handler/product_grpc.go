package handler

import (
	"context"

	"ecommerce/internal/domain"
	pb "ecommerce/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ProductGrpcHandler adalah manajer khusus untuk jalur gRPC
type ProductGrpcHandler struct {
	pb.UnimplementedProductServiceServer
	productUsecase domain.ProductUsecase
}

// NewProductGrpcHandler adalah constructor-nya
func NewProductGrpcHandler(pu domain.ProductUsecase) *ProductGrpcHandler {
	return &ProductGrpcHandler{
		productUsecase: pu,
	}
}

// GetProductByID adalah implementasi dari kontrak RPC yang kita tulis di file .proto
func (h *ProductGrpcHandler) GetProductByID(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	// 1. Panggil Usecase
	product, err := h.productUsecase.GetProductByID(ctx, int(req.GetId()))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Produk tidak ditemukan: %v", err)
	}

	// 2. Konversi tipe decimal shopspring ke float64
	price, _ := product.Price.Float64()

	// 3. Masukkan data dari database ke dalam "Amplop Protobuf"
	return &pb.GetProductResponse{
		Id:    int32(product.ID),
		Name:  product.Name,
		Price: price,
		Stock: int32(product.Stock),
	}, nil
}
