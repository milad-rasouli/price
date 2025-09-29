package routes

import (
	"github.com/gin-gonic/gin"
)

type Router interface {
	SetupRoutes(engine *gin.Engine)
}

func CreateRouters(
	priceRouter *PriceRouter,
	cron *CronRouter,
	healthRouter *HealthRouter,
) []Router {
	return []Router{
		healthRouter,
		priceRouter,
		cron,
	}
}
