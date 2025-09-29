package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/milad-rasouli/price/internal/app/api/controllers"
)

type PriceRouter struct {
	priceController *controller.PriceController
}

func NewPriceRouter(priceController *controller.PriceController) *PriceRouter {
	return &PriceRouter{priceController: priceController}
}

func (pr *PriceRouter) SetupRoutes(router *gin.Engine) {
	g := router.Group("/prices")
	{
		g.GET("/history", pr.priceController.GetHistory)
		g.GET("/latest", pr.priceController.GetLatest)
	}
}
