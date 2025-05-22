package api

import (
	"backend/config"
	"backend/routes"
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
	s.router = routes.SetupRouter(s.db, s.config)
}

func (s *Server) Start() error {
	return s.router.Run(":" + s.config.Server.Port)
} 