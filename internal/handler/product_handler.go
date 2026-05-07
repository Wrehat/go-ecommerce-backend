package handler

import (
	"ecommerce/internal/domain"
	"ecommerce/internal/usecase"
	"ecommerce/pkg/response"
	"errors"
	"net/http"
	"strconv"
	"time"

	// "strconv"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type ProductRequest struct {
	SKU   string          `json:"sku" binding:"required"`
	Name  string          `json:"name" binding:"required,min=3"`
	Price decimal.Decimal `json:"price" binding:"required"`
	Stock int             `json:"stock" binding:"required,gte=0"`
}

type ProductResponse struct {
	ID        int             `json:"id"`
	SKU       string          `json:"sku"`
	Name      string          `json:"name"`
	Price     decimal.Decimal `json:"price"`
	Stock     int             `json:"stock"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

func ToProductResponse(p domain.Product) ProductResponse {
	return ProductResponse{
		ID:        p.ID,
		SKU:       p.SKU,
		Name:      p.Name,
		Price:     p.Price,
		Stock:     p.Stock,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

type ProductHandler struct {
	// Menyimpan referensi ke usecase yang akan digunakan untuk operasi bisnis
	usecase *usecase.ProductUsecase
}

// Konstruktor untuk membuat instance baru dari ProductHandler
func NewProductHandler(u *usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{
		usecase: u,
	}
}

// ==========================================
// FUNGSI-FUNGSI ENDPOINT
// ==========================================
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	// wadah request
	var reqPayload ProductRequest

	// Validasi request
	if err := c.ShouldBindJSON(&reqPayload); err != nil {
		response.JSON(c, http.StatusBadRequest, "Data product tidak valid", nil)
		return
	}

	if reqPayload.Price.IsNegative() || reqPayload.Price.IsZero() {
		response.JSON(c, http.StatusBadRequest, "Harga harus lebih besar dari nol", nil)
		return
	}

	// Buat product
	newProduct := domain.Product{
		SKU:   reqPayload.SKU,
		Name:  reqPayload.Name,
		Price: reqPayload.Price,
		Stock: reqPayload.Stock,
	}

	// Request Context
	ctx := c.Request.Context()

	// Panggil usecase
	product, err := h.usecase.CreateProduct(ctx, newProduct)
	if err != nil {
		if errors.Is(err, domain.ErrSKUDuplicate) {
			response.JSON(c, http.StatusConflict, err.Error(), nil)
			return
		}
		response.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// return response
	response.JSON(c, http.StatusCreated, "Product berhasil dibuat", ToProductResponse(product))
}

func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	// Request Context
	ctx := c.Request.Context()

	products, err := h.usecase.GetAllProducts(ctx)
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	productResponse := make([]ProductResponse, 0)
	for _, p := range products {
		productResponse = append(productResponse, ToProductResponse(p))

	}

	response.JSON(c, http.StatusOK, "Daftar produk", productResponse)
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
	// Todo : Ambil id dan context
	id := c.Param("id")
	idReq, err := strconv.Atoi(id)
	if err != nil {
		response.JSON(c, http.StatusBadRequest, "ID harus berupa angka", nil)
		return
	}
	if idReq <= 0 {
		response.JSON(c, http.StatusBadRequest, "ID harus berupa angka positif", nil)
		return
	}

	ctx := c.Request.Context()

	// Todo : Panggil usecase
	product, err := h.usecase.GetProductByID(ctx, idReq)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			response.JSON(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		response.JSON(c, http.StatusInternalServerError, "Terjadi kesalahan pada server", nil)
		return
	}

	// Todo : Buat response
	response.JSON(c, http.StatusOK, "Berhasil mengambil data product", ToProductResponse(product))

}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	// Todo : ambil context
	ctx := c.Request.Context()
	// Todo : ambil & validasi ID
	id := c.Param("id")
	idReq, err := strconv.Atoi(id)
	if err != nil {
		response.JSON(c, http.StatusBadRequest, "Id Product tidak valid", nil)
		return
	}
	// Todo : buat wadah & bind request
	var reqPayload ProductRequest
	if err := c.ShouldBindJSON(&reqPayload); err != nil {
		response.JSON(c, http.StatusBadRequest, "Format data product tidak valid", nil)
		return
	}

	if reqPayload.Price.IsNegative() || reqPayload.Price.IsZero() {
		response.JSON(c, http.StatusBadRequest, "Harga harus lebih besar dari nol", nil)
		return
	}

	reqProduct := domain.Product{
		SKU:   reqPayload.SKU,
		Name:  reqPayload.Name,
		Price: reqPayload.Price,
		Stock: reqPayload.Stock,
	}

	// panggil usecase dan handle error
	updatedProduct, err := h.usecase.UpdateProduct(ctx, idReq, reqProduct)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			response.JSON(c, http.StatusNotFound, "Data product tidak ditemukan", nil)
			return
		}
		if errors.Is(err, domain.ErrSKUDuplicate) {
			response.JSON(c, http.StatusConflict, "SKU payload harus unik", nil)
			return
		}
		response.JSON(c, http.StatusInternalServerError, "Terjadi kesalahan pada server", nil)
		return
	}

	// return response
	response.JSON(c, http.StatusOK, "Product berhasil diperbarui", ToProductResponse(updatedProduct))
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	// Todo: ambil id dan context
	id := c.Param("id")
	idReq, err := strconv.Atoi(id)
	if err != nil {
		response.JSON(c, http.StatusBadRequest, "Id product tidak valid", nil)
		return
	}

	ctx := c.Request.Context()

	// Todo: panggil usecase & handle error
	if err := h.usecase.DeleteProduct(ctx, idReq); err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			response.JSON(c, http.StatusNotFound, "Data product tidak ditemukan", nil)
			return
		}
		response.JSON(c, http.StatusInternalServerError, "Terjadi kesalahan pada server", nil)
		return
	}

	// Todo: return success response
	response.JSON(c, http.StatusOK, "Data berhasil dihapus", nil)
}

// CreateProduct godoc
// @Summary Tambah Produk Baru
// @Description Menyimpan produk baru ke dalam sistem
// @Tags Products
// @Accept json
// @Produce json
// @Param request body ProductRequest true "Format JSON untuk produk baru"
// @Success 201 {object} APIResponse
// @Failure 400 {object} APIResponse
// @Router /products [post]
// func (h *ProductHandler) CreateProduct(c *gin.Context) {
// 	// Binding JSON ke struct Product
// 	var req ProductRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		respondJSON(c, http.StatusBadRequest, "Format JSON tidak valid", nil)
// 		return
// 	}

// 	newProduct := domain.Product{
// 		Name:  req.Name,
// 		Price: req.Price,
// 		Stock: req.Stock,
// 	}

// 	createdProduct, err := h.usecase.CreateProduct(newProduct)
// 	if err != nil {
// 		respondJSON(c, http.StatusBadRequest, err.Error(), nil)
// 		return
// 	}

// 	respondJSON(c, http.StatusCreated, "Produk berhasil dibuat", createdProduct)

// }

// GET /products/:id
// GetProductByID godoc
// @Summary Dapatkan Produk berdasarkan ID
// @Description Mengambil informasi produk berdasarkan ID yang diberikan
// @Tags Products
// @Accept json
// @Produce json
// @Param id path int true "ID Produk"
// @Success 200 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Router /products/{id} [get]
// func (h *ProductHandler) GetProductByID(c *gin.Context) {
// 	productID := c.Param("id")
// 	productIDInt, err := strconv.Atoi(productID)
// 	if err != nil {
// 		respondJSON(c, http.StatusBadRequest, "ID harus berupa angka", nil)
// 		return
// 	}

// 	if productIDInt <= 0 {
// 		respondJSON(c, http.StatusBadRequest, "ID produk tidak valid", nil)
// 		return
// 	}

// 	product, err := h.usecase.GetProductByID(productIDInt)
// 	if err != nil {
// 		respondJSON(c, http.StatusNotFound, err.Error(), nil)
// 		return
// 	}

// 	respondJSON(c, http.StatusOK, "Produk ditemukan", product)
// }

// GET /products
// GetAllProducts godoc
// @Summary Dapatkan Semua Produk
// @Description Mengambil daftar semua produk yang tersedia dalam sistem
// @Tags Products
// @Accept json
// @Produce json
// @Success 200 {object} APIResponse
// @Router /products [get]

// PUT /products/:id
// UpdateProduct godoc
// @Summary Perbarui Produk
// @Description Memperbarui informasi produk berdasarkan ID yang diberikan
// @Tags Products
// @Accept json
// @Produce json
// @Param id path int true "ID Produk"
// @Param request body ProductRequest true "Format JSON untuk pembaruan produk"
// @Success 200 {object} APIResponse
// @Failure 400 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Router /products/{id} [put]
// func (h *ProductHandler) UpdateProduct(c *gin.Context) {
// 	productID := c.Param("id")
// 	productIDInt, err := strconv.Atoi(productID)
// 	if err != nil {
// 		respondJSON(c, http.StatusBadRequest, "ID harus berupa angka", nil)
// 		return
// 	}

// 	var req ProductRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		respondJSON(c, http.StatusBadRequest, "Format JSON tidak valid", nil)
// 		return
// 	}

// 	newUpdateProduct := domain.Product{
// 		Name:  req.Name,
// 		Price: req.Price,
// 		Stock: req.Stock,
// 	}

// 	updatedProduct, err := h.usecase.UpdateProduct(productIDInt, newUpdateProduct)
// 	if err != nil {
// 		respondJSON(c, http.StatusBadRequest, err.Error(), nil)
// 		return
// 	}

// 	respondJSON(c, http.StatusOK, "Produk berhasil diperbarui", updatedProduct)
// }

// DELETE /products/:id
// DeleteProduct godoc
// @Summary Hapus Produk
// @Description Menghapus produk berdasarkan ID yang diberikan
// @Tags Products
// @Accept json
// @Produce json
// @Param id path int true "ID Produk"
// @Success 200 {object} APIResponse
// @Failure 400 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Router /products/{id} [delete]
// func (h *ProductHandler) DeleteProduct(c *gin.Context) {
// 	productID := c.Param("id")
// 	productIDInt, err := strconv.Atoi(productID)
// 	if err != nil {
// 		respondJSON(c, http.StatusBadRequest, "ID harus berupa angka", nil)
// 		return
// 	}

// 	err = h.usecase.DeleteProduct(productIDInt)
// 	if err != nil {
// 		respondJSON(c, http.StatusNotFound, err.Error(), nil)
// 		return
// 	}

// 	respondJSON(c, http.StatusOK, "Produk berhasil dihapus", nil)
// }
