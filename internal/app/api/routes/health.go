package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/milad-rasouli/price/internal/app/api/controllers"
)

type HealthRouter struct {
	healthController *controller.HealthController
}

func NewHealthRouter(healthController *controller.HealthController) *HealthRouter {
	return &HealthRouter{healthController: healthController}
}

func (rh *HealthRouter) SetupRoutes(router *gin.Engine) {
	g := router.Group("/")
	{
		g.GET("/liveness", rh.healthController.Liveness)
		g.GET("/", rh.healthController.Liveness)
		g.GET("/readiness", rh.healthController.Readiness)
	}
}
