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
	"ecommerce/pkg/logger"
	"ecommerce/pkg/telemetry"
	pb "ecommerce/proto"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
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

	logger.InitLogger(cfg.AppEnv)
	defer logger.Sync()
	logger.Log.Info("Sistem E-Commerce mulai dijalankan...", zap.String("env", cfg.AppEnv))

	tp, err := telemetry.InitTracer("ecommerce-api")
	if err != nil {
		logger.Log.Error("Gagal inisialisasi OTel Tracer", zap.Error(err))
		panic("Gagal menyalakan OpenTelemetry")
	}

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			logger.Log.Error("Error saat mematikan trace", zap.Error(err))
		}
	}()

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

	productGrpcHandler := handler.NewProductGrpcHandler(productUsecase)

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
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(otelgin.Middleware("ecommerce-api"))
	pprof.Register(r)
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger())
	r.Use(middleware.RateLimiterTokenBucket(redisClient, 1, 10))
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

	setupRoutes()

	// =============================================================
	// 4. SERVER RUNNER & LIFECYCLE MANAGEMENT (Menjalankan Mesin)
	// =============================================================

	grpcServer := grpc.NewServer()
	pb.RegisterProductServiceServer(grpcServer, productGrpcHandler)

	grpcListener, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Log.Error("Gagal membuka port gRPC", zap.Error(err))
		panic(err)
	}

	go func() {
		logger.Log.Info("🚀 gRPC Server berjalan di jalur biner", zap.String("port", "50051"))
		if err := grpcServer.Serve(grpcListener); err != nil {
			logger.Log.Error("Error saat menjalankan gRPC server", zap.Error(err))
		}
	}()

	srv := &http.Server{
		Addr:         ":" + cfg.AppPort,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// Jalankan server di background goroutine agar tidak memblokir shutdown listener
	go func() {
		logger.Log.Info(
			"Api Server berjalan",
			zap.String("port", cfg.AppPort),
			zap.String("docs", "http://localhost:"+cfg.AppPort+"/swagger/index.html"),
		)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Error("Error saat menjalankan server", zap.Error(err))
		}
	}()

	// Dengarkan sinyal OS untuk mematikan server secara aman
	handleGracefulShutdown(srv, grpcServer)
}

// =============================================================
// EXTRACTION FUNCTIONS (Fungsi Pembantu di Luar Main)
// =============================================================

// handleGracefulShutdown mengisolasi logika mekanik shutdown agar fungsi main() tetap bersih
func handleGracefulShutdown(srv *http.Server, grpcServer *grpc.Server) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Menunggu tombol Ctrl+C atau perintah kill di terminal
	<-ctx.Done()

	logger.Log.Info("Sinyal mati diterima. Mematikan seluruh mesin...")

	// Matikan gRPC Server
	grpcServer.GracefulStop()
	logger.Log.Info("✅ gRPC Server dimatikan dengan aman.")

	// Matikan HTTP Server
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Log.Error("Error saat mematikan HTTP server", zap.Error(err))
	}

	logger.Log.Info("✅ HTTP Server dimatikan dengan aman.")
}
