package controller

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/milad-rasouli/price/internal/app/api/response"
	"github.com/milad-rasouli/price/internal/infrastructure/postgresql"
)

type HealthController struct {
	lastReady time.Time
	logger    *slog.Logger
	pg        *postgresql.Postgres
}

func NewHealthController(logger *slog.Logger, pg *postgresql.Postgres) *HealthController {
	return &HealthController{
		lastReady: time.Now(),
		pg:        pg,
		logger:    logger.With("layer", "HealthController"),
	}
}

// Liveness godoc
// @Summary Liveness probe
// @Description Used by Kubernetes or monitoring tools to check if the service is alive.
// @Tags health
// @Produce json
// @Success 200 {object} response.Response[any] "Service is alive"
// @Failure 503 {object} response.Response[any] "Service unavailable"
// @Router /liveness [get]
func (hc *HealthController) Liveness(c *gin.Context) {
	lg := hc.logger.With("method", "Liveness")

	if hc.lastReady.Compare(time.Now().Add(-(time.Minute * 5))) == -1 {
		lg.Warn("5 minutes of Readiness failure, making liveness fail to restart...")
		response.Custom(c, http.StatusServiceUnavailable, nil, "")
		return
	}

	response.Ok(c, nil, "")
}

// Readiness godoc
// @Summary Readiness probe
// @Description Verifies if dependencies (e.g., PostgreSQL) are healthy and service can handle requests.
// @Tags health
// @Produce json
// @Success 200 {object} response.Response[any] "Service is ready"
// @Failure 503 {object} response.Response[any] "Service not ready"
// @Router /readiness [get]
func (hc *HealthController) Readiness(c *gin.Context) {
	lg := hc.logger.With("method", "Readiness")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if err := hc.pg.HealthCheck(ctx); err != nil {
		lg.Error("pgx health check failed", "error", err)
		response.Custom(c, http.StatusServiceUnavailable, nil, "")
		return
	}

	hc.lastReady = time.Now()
	response.Ok(c, nil, "")
}
