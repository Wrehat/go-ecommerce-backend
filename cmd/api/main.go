package main

import (
	"context"
	_ "ecommerce/docs"
	"ecommerce/internal/config"
	"ecommerce/internal/database"
	"ecommerce/internal/handler"
	"ecommerce/internal/repository"
	"ecommerce/internal/usecase"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Toko Arsitek API
// @version 1.0
// @description Ini adalah dokumentasi API untuk E-Commerce Kelas Enterprise.
// @host localhost:8080
// @BasePath /api/v1
func main() {
	cfg := config.LoadConfig()

	// Koneksi DB
	dbPool := database.ConnectDB(cfg.DBUri)
	defer dbPool.Close()
	database.RunMigrations(cfg.DBUri)

	// Koneksi Redis
	redisClient := database.ConnRedis(cfg.RedisURL)
	defer redisClient.Close()

	productRepo := repository.NewProductRepository(dbPool)

	productUsecase := usecase.NewProductUsecase(productRepo, redisClient)

	productHandler := handler.NewProductHandler(productUsecase)

	orderRepo := repository.NewOrderRepository(dbPool)

	orderUsecase := usecase.NewOrderUsecase(orderRepo, productRepo)

	orderHandler := handler.NewOrderHandler(orderUsecase)

	// Inisialisasi Gin Router
	router := gin.Default()

	v1 := router.Group("/api/v1")

	productGroup := v1.Group("/products")
	{
		productGroup.GET("/", productHandler.GetAllProducts)
		productGroup.GET("/:id", productHandler.GetProductByID)
	}

	orderGroup := v1.Group("/orders")
	{
		orderGroup.POST("/checkout", orderHandler.CheckOut)
	}

	adminGroup := v1.Group("/admin")
	{
		adminGroup.POST("/products", productHandler.CreateProduct)
		adminGroup.PUT("/products/:id", productHandler.UpdateProduct)
		adminGroup.DELETE("/products/:id", productHandler.DeleteProduct)
	}

	// Buat endpoint sederhana(Ping)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "PONG! Server Gin berjalan dengan sempurna!",
		})
	})
	// Rute untuk melihat Buku Panduan Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	srv := &http.Server{
		Addr:         ":" + cfg.AppPort,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	go func() {
		fmt.Printf("🚀 API Server berjalan di http://localhost:%v\n", cfg.AppPort)
		fmt.Println("📖 Dokumentasi API tersedia di http://localhost:" + cfg.AppPort + "/swagger/index.html")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Error saat menjalankan server %v\n", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		fmt.Printf("Error saat mematikan server %v\n", err)
	}

	fmt.Println("Server dimatikan")
}
