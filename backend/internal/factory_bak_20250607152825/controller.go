package factory

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller struct {
	service *Service
}

func NewController(service *Service) *Controller {
	return &Controller{service: service}
}

// GetFactoryList 获取工厂清单
func (c *Controller) GetFactoryList(ctx *gin.Context) {
	factories, err := c.service.GetFactories()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, factories)
} 