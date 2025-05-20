package server

import (
	"api-gateway/config"
	handler "api-gateway/internal/handlers"
	"api-gateway/internal/middleware"
	"fmt"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Server encapsulates the Gin router
type Server struct {
	router *gin.Engine
	cfg    *config.Config
}

// NewServer initializes the server with routes and middleware
func NewServer(cfg *config.Config) *Server {
	r := gin.New()

	// ✅ Custom CORS config that allows Authorization headers and local frontend
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:5500"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Initialize handlers
	h, err := handler.NewHandler(cfg)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize handler: %v", err))
	}

	// Apply global middleware
	r.Use(middleware.LoggingMiddleware())
	r.Use(middleware.TelemetryMiddleware())

	// Define routes
	api := r.Group("/api")

	// Unprotected routes (no authentication required)
	api.POST("/users/register", h.RegisterUser)
	api.POST("/users/authenticate", h.AuthenticateUser)

	// Protected routes (require authentication)
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(cfg))
	{
		// Inventory routes
		protected.POST("/inventory", h.CreateProduct)
		protected.GET("/inventory/:id", h.GetProduct)
		protected.PUT("/inventory/:id", h.UpdateProduct)
		protected.DELETE("/inventory/:id", h.DeleteProduct)
		protected.GET("/inventory", h.ListProducts)

		// Order routes
		protected.POST("/orders", h.CreateOrder)
		protected.GET("/orders/:id", h.GetOrder)
		protected.PUT("/orders/:id/status", h.UpdateOrderStatus)
		protected.GET("/orders", h.ListUserOrders)

		// User routes
		protected.GET("/users/:id", h.GetUserProfile)
	}

	return &Server{
		router: r,
		cfg:    cfg,
	}
}

// Start runs the server
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%s", s.cfg.Port)
	return s.router.Run(addr)
}
