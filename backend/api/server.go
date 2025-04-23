package api

import (
	"sewingmast-backend/config"
	"sewingmast-backend/routes"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	config *config.Config
	router *gin.Engine
	db     *gorm.DB
}

func NewServer(cfg *config.Config, db *gorm.DB) *Server {
	server := &Server{
		config: cfg,
		router: gin.Default(),
		db:     db,
	}

	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	// Register auth routes
	routes.RegisterAuthRoutes(s.router, s.db)

	// Protected routes
	protected := s.router.Group("/api")
	// protected.Use(middleware.JWTAuth())
	{
		// Designer routes
		designer := protected.Group("/designer")
		{
			designer.GET("/orders", s.handleGetDesignerOrders)
			designer.POST("/orders", s.handleCreateOrder)
		}

		// Factory routes
		factory := protected.Group("/factory")
		{
			factory.GET("/orders", s.handleGetFactoryOrders)
			factory.PUT("/orders/:id", s.handleUpdateOrderStatus)
		}
	}
}

func (s *Server) Start() error {
	return s.router.Run(":" + s.config.Server.Port)
} 