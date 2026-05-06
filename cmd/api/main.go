package main

import (
	"context"
	_ "ecommerce/docs"
	"ecommerce/internal/database"
	"ecommerce/internal/handler"
	"ecommerce/internal/repository"
	"ecommerce/internal/usecase"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Toko Arsitek API
// @version 1.0
// @description Ini adalah dokumentasi API untuk E-Commerce Kelas Enterprise.
// @host localhost:8080
// @BasePath /api/v1
func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv("DB_URI")
	if dsn == "" {
		fmt.Println("DB_URI env var belum di set")
		os.Exit(1)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbPool := database.ConnectDB(dsn)
	defer dbPool.Close()
	database.RunMigrations(dsn)


	productRepo := repository.NewProductRepository(dbPool)

	productUsecase := usecase.NewProductUsecase(productRepo)

	productHandler := handler.NewProductHandler(productUsecase)

	// Inisialisasi Gin Router
	router := gin.Default()

	v1 := router.Group("/api/v1")

	productGroup := v1.Group("/products")
	{
		productGroup.POST("/", productHandler.CreateProduct)
		productGroup.GET("/", productHandler.GetAllProducts)
		// productGroup.GET("/:id", productHandler.GetProductByID)
		// productGroup.PUT("/:id", productHandler.UpdateProduct)
		// productGroup.DELETE("/:id", productHandler.DeleteProduct)
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
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	go func() {
		fmt.Printf("🚀 API Server berjalan di http://localhost:%v\n", port)
		fmt.Println("📖 Dokumentasi API tersedia di http://localhost:8080/swagger/index.html")
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
