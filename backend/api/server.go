package api

import (
	"sewingmast-backend/config"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config *config.Config
	router *gin.Engine
}

func NewServer(cfg *config.Config) *Server {
	server := &Server{
		config: cfg,
		router: gin.Default(),
	}

	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	// Public routes
	public := s.router.Group("/api")
	{
		public.POST("/login", s.handleLogin)
		public.POST("/register", s.handleRegister)
	}

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