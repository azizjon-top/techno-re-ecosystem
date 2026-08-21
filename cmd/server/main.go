package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"techno-re-ecosystem/internal/config"
	"techno-re-ecosystem/internal/logger"
	"techno-re-ecosystem/internal/middleware"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	if err := logger.Init(cfg.Server.Environment); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Get().Close()

	logger.Info("Starting Techno RE Ecosystem Server",
		"environment", cfg.Server.Environment,
		"port", cfg.Server.Port,
	)

	// Initialize Gin router
	if cfg.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()

	// Apply middleware
	router.Use(middleware.ErrorHandler())
	router.Use(middleware.RequestLogger())
	router.Use(middleware.CORS())
	router.Use(middleware.ContentTypeJSON())
	router.Use(middleware.RequestTimeout(cfg.Server.ReadTimeout))

	// Health check endpoint
	router.GET("/health", healthCheck)

	// API v1 routes (placeholder for future implementation)
	v1 := router.Group("/api/v1")
	{
		// Auth routes
		v1.POST("/auth/register", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "register endpoint"})
		})
		v1.POST("/auth/login", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "login endpoint"})
		})

		// User routes
		v1.GET("/users/:id", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "get user endpoint"})
		})

		// Product routes
		v1.GET("/products", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "list products endpoint"})
		})
		v1.POST("/products", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "create product endpoint"})
		})

		// Order routes
		v1.GET("/orders", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "list orders endpoint"})
		})
		v1.POST("/orders", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "create order endpoint"})
		})

		// Wallet routes
		v1.GET("/wallet/balance", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "get wallet balance endpoint"})
		})
		v1.POST("/wallet/transfer", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "transfer tokens endpoint"})
		})

		// Chat routes
		v1.GET("/chats", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "list chats endpoint"})
		})
		v1.POST("/chats/:id/messages", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "send message endpoint"})
		})

		// Video routes
		v1.GET("/videos", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "list videos endpoint"})
		})
		v1.POST("/videos", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "upload video endpoint"})
		})

		// Mining routes
		v1.POST("/mining/start", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "start mining session endpoint"})
		})
		v1.POST("/mining/validate", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "validate fact endpoint"})
		})

		// Campaign routes
		v1.GET("/campaigns", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "list campaigns endpoint"})
		})
		v1.POST("/campaigns", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "create campaign endpoint"})
		})
	}

	// 404 handler
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    "NOT_FOUND",
			"message": "Endpoint not found",
		})
	})

	// Create server
	srv := &http.Server{
		Addr:         cfg.Server.Host + ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Start server in goroutine
	go func() {
		logger.Info("Server starting", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed to start", "error", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	logger.Info("Server shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown error", "error", err)
		os.Exit(1)
	}

	logger.Info("Server stopped successfully")
}

// healthCheck returns health status of the server
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"timestamp": time.Now(),
	})
}
