package main

import (
	"context"
	_ "ecommerce/docs"
	"ecommerce/internal/config"
	"ecommerce/internal/database"
	"ecommerce/internal/handler"
	"ecommerce/internal/middleware"
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
// @description Ini adalah dokumentasi API untuk E-Commerce.
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// =============================================================
	// 1. INFRASTRUCTURE & CONFIG INITIALIZATION (Koneksi & Env)
	// =============================================================
	cfg := config.LoadConfig()
	dbPool := database.ConnectDB(cfg.DBUri)
	defer dbPool.Close()
	database.RunMigrations(cfg.DBUri)
	redisClient := database.ConnRedis(cfg.RedisURL)
	defer redisClient.Close()

	// =============================================================
	// 2. DEPENDENCY INJECTION REGISTRY (Penyambungan Komponen)
	// =============================================================

	// Product Module
	productRepo := repository.NewProductRepository(dbPool)
	productUsecase := usecase.NewProductUsecase(productRepo, redisClient)
	productHandler := handler.NewProductHandler(productUsecase)

	// Order Module
	orderRepo := repository.NewOrderRepository(dbPool)
	orderUsecase := usecase.NewOrderUsecase(orderRepo, productRepo)
	orderHandler := handler.NewOrderHandler(orderUsecase)

	// User / Auth Module
	userRepo := repository.NewUserRepository(dbPool)
	userUsecase := usecase.NewUserUsecase(userRepo, cfg.JWTSecret)
	userHandler := handler.NewUserHandler(userUsecase)

	// =============================================================
	// 3. ROUTER & API MAPPING (Pemetaan Rute HTTP)
	// =============================================================
	// r := gin.Default()
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger())
	r.Use(middleware.RateLimiter(5, 10))
	authMiddleware := middleware.RequireAuth(cfg.JWTSecret)
	setupRoutes := func() {
		v1 := r.Group("/api/v1")

		// Public Product Routes
		productGroup := v1.Group("/products")
		{
			productGroup.GET("", productHandler.GetAllProducts)
			productGroup.GET("/:id", productHandler.GetProductByID)
		}

		// Order Routes
		orderGroup := v1.Group("/orders").Use(authMiddleware)
		{
			orderGroup.POST("/checkout", orderHandler.CheckOut)
		}

		// Admin Routes
		adminGroup := v1.Group("/admin").Use(authMiddleware, middleware.RequireRole("admin"))
		{
			adminGroup.POST("/products", productHandler.CreateProduct)
			adminGroup.PUT("/products/:id", productHandler.UpdateProduct)
			adminGroup.DELETE("/products/:id", productHandler.DeleteProduct)
		}

		// User / Auth Routes
		userGroup := v1.Group("/users")
		{
			userGroup.POST("/register", userHandler.Register)
			userGroup.POST("/login", userHandler.Login)
		}

		// Base Utility Routes
		r.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "PONG! Server Gin berjalan dengan sempurna!",
			})
		})
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// Jalankan pemetaan rute
	setupRoutes()

	// =============================================================
	// 4. SERVER RUNNER & LIFECYCLE MANAGEMENT (Menjalankan Mesin)
	// =============================================================
	srv := &http.Server{
		Addr:         ":" + cfg.AppPort,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// Jalankan server di background goroutine agar tidak memblokir shutdown listener
	go func() {
		fmt.Printf("🚀 API Server berjalan di http://localhost:%v\n", cfg.AppPort)
		fmt.Println("📖 Dokumentasi API tersedia di http://localhost:" + cfg.AppPort + "/swagger/index.html")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("❌ Error saat menjalankan server: %v\n", err)
		}
	}()

	// Dengarkan sinyal OS untuk mematikan server secara aman
	handleGracefulShutdown(srv)
}

// =============================================================
// EXTRACTION FUNCTIONS (Fungsi Pembantu di Luar Main)
// =============================================================

// handleGracefulShutdown mengisolasi logika mekanik shutdown agar fungsi main() tetap bersih
func handleGracefulShutdown(srv *http.Server) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Menunggu tombol Ctrl+C atau perintah kill di terminal
	<-ctx.Done()

	fmt.Println("\n⏳ Sinyal mati diterima. Mematikan server secara anggun (Graceful Shutdown)...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		fmt.Printf("❌ Error saat mematikan server: %v\n", err)
	}

	fmt.Println("🛑 Server berhasil dimatikan dengan aman sepenuhnya.")
}
