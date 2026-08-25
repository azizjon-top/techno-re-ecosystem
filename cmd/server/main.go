package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/azizjon-top/techno-re-ecosystem/internal/config"
	"github.com/azizjon-top/techno-re-ecosystem/internal/handler"
	"github.com/azizjon-top/techno-re-ecosystem/internal/logger"
	"github.com/azizjon-top/techno-re-ecosystem/internal/middleware"
	"github.com/azizjon-top/techno-re-ecosystem/internal/repository"
	"github.com/azizjon-top/techno-re-ecosystem/internal/service"
)

func main() {
	cfg := config.Load()

	if err := logger.Init(cfg.Server.Environment); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Get().Close()

	logger.Info("Starting Techno RE Ecosystem Server",
		zap.String("environment", cfg.Server.Environment),
		zap.String("port", cfg.Server.Port),
	)

	// ---- Repositories (in-memory; swap for DB implementations later) ----
	userRepo := repository.NewInMemoryUserRepository()
	productRepo := repository.NewInMemoryProductRepository()
	orderRepo := repository.NewInMemoryOrderRepository()
	walletRepo := repository.NewInMemoryWalletRepository()

	// ---- Services ----
	jwtSvc := service.NewJWTService(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenTTL,
		cfg.JWT.RefreshTokenTTL,
		cfg.JWT.TokenIssuer,
		cfg.JWT.TokenAudience,
	)
	userSvc := service.NewUserService(userRepo, walletRepo, jwtSvc)
	productSvc := service.NewProductService(productRepo, orderRepo, walletRepo)
	walletSvc := service.NewWalletService(walletRepo)

	// ---- Handlers ----
	authHandler := handler.NewAuthHandler(userSvc)
	userHandler := handler.NewUserHandler(userSvc)
	productHandler := handler.NewProductHandler(productSvc)
	walletHandler := handler.NewWalletHandler(walletSvc)

	// ---- Router ----
	if cfg.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()
	router.Use(middleware.ErrorHandler())
	router.Use(middleware.RequestLogger())
	router.Use(middleware.CORS())
	router.Use(middleware.ContentTypeJSON())
	router.Use(middleware.RequestTimeout(cfg.Server.ReadTimeout))

	router.GET("/health", healthCheck)

	v1 := router.Group("/api/v1")

	// Public auth routes
	auth := v1.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// Protected routes
	protected := v1.Group("")
	protected.Use(middleware.Auth(jwtSvc))
	{
		// Users
		protected.GET("/users/me", userHandler.GetMe)
		protected.GET("/users/:id", userHandler.GetUser)

		// Products
		protected.GET("/products", productHandler.ListProducts)
		protected.POST("/products", middleware.RequireRole("seller", "admin"), productHandler.CreateProduct)

		// Orders
		protected.GET("/orders", productHandler.ListOrders)
		protected.POST("/orders", productHandler.CreateOrder)

		// Wallet
		protected.GET("/wallet/balance", walletHandler.GetBalance)
		protected.POST("/wallet/transfer", walletHandler.Transfer)

		// Chat (stubs — real-time WebSocket to be added later)
		protected.GET("/chats", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"chats": []interface{}{}})
		})
		protected.POST("/chats/:id/messages", func(c *gin.Context) {
			c.JSON(http.StatusNotImplemented, gin.H{"message": "WebSocket chat coming soon"})
		})

		// Videos (stubs)
		protected.GET("/videos", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"videos": []interface{}{}})
		})
		protected.POST("/videos", func(c *gin.Context) {
			c.JSON(http.StatusNotImplemented, gin.H{"message": "video upload coming soon"})
		})

		// Mining (stubs)
		protected.POST("/mining/start", func(c *gin.Context) {
			c.JSON(http.StatusNotImplemented, gin.H{"message": "mining coming soon"})
		})
		protected.POST("/mining/validate", func(c *gin.Context) {
			c.JSON(http.StatusNotImplemented, gin.H{"message": "mining coming soon"})
		})

		// Campaigns (stubs)
		protected.GET("/campaigns", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"campaigns": []interface{}{}})
		})
		protected.POST("/campaigns", func(c *gin.Context) {
			c.JSON(http.StatusNotImplemented, gin.H{"message": "campaigns coming soon"})
		})
	}

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    "NOT_FOUND",
			"message": "Endpoint not found",
		})
	})

	// ---- HTTP server ----
	srv := &http.Server{
		Addr:         cfg.Server.Host + ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		logger.Info("Server starting", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Server shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown error", zap.Error(err))
		os.Exit(1)
	}
	logger.Info("Server stopped successfully")
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now(),
	})
}
