package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"api-gateway/config"
	"api-gateway/internal/clients"
	"api-gateway/internal/handlers"
	"api-gateway/internal/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	inventoryClient, err := clients.NewInventoryClient(cfg.InventoryAddr)
	if err != nil {
		log.Fatalf("Failed to connect to Inventory Service: %v", err)
	}
	defer inventoryClient.Close()

	orderClient, err := clients.NewOrderClient(cfg.OrderAddr)
	if err != nil {
		log.Fatalf("Failed to connect to Order Service: %v", err)
	}
	defer orderClient.Close()

	router := gin.New()

	router.Use(
		middleware.LoggingMiddleware(logger),
		gin.Recovery(),
	)
	authHandler := handlers.NewAuthHandler(cfg)
	router.POST("/auth/login", authHandler.Login) 

	api := router.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(cfg))
		{
			inventoryHandler := handlers.NewInventoryHandler(inventoryClient)
			protected.POST("/products", inventoryHandler.CreateProduct)
			protected.GET("/products/:id", inventoryHandler.GetProduct)
			protected.PUT("/products/:id", inventoryHandler.UpdateProduct)
			protected.DELETE("/products/:id", inventoryHandler.DeleteProduct)
			protected.GET("/products", inventoryHandler.ListProducts)

			orderHandler := handlers.NewOrderHandler(orderClient)
			protected.POST("/orders", orderHandler.CreateOrder)
			protected.GET("/orders/:id", orderHandler.GetOrder)
			protected.PUT("/orders/:id/status", orderHandler.UpdateOrderStatus)
			protected.GET("/orders", orderHandler.ListUserOrders)
		}
	}

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		logger.Info("Starting API Gateway", zap.String("port", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")
	if err := srv.Close(); err != nil {
		logger.Error("Server shutdown error", zap.Error(err))
	}
	logger.Info("Server exited")
}
