package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/milad-rasouli/price/internal/app/api/controllers"
)

type CronRouter struct {
	cronController *controller.CronController
}

func NewCronRouter(cronController *controller.CronController) *CronRouter {
	return &CronRouter{cronController: cronController}
}

func (cr *CronRouter) SetupRoutes(router *gin.Engine) {
	g := router.Group("/cron")
	{
		g.GET("/update-prices", cr.cronController.UpdatePrice)
	}
}
