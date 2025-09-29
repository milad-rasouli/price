package controller

import (
	"context"
	"errors"
	"github.com/milad-rasouli/price/internal/providers/currency"
	"github.com/milad-rasouli/price/internal/service"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/milad-rasouli/price/internal/app/api/response"
)

type CronController struct {
	logger  *slog.Logger
	service service.PriceService
}

func NewCronController(logger *slog.Logger, svc service.PriceService) *CronController {
	return &CronController{
		logger:  logger.With("layer", "CronController"),
		service: svc,
	}
}

func (pc *CronController) UpdatePrice(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	err := pc.service.InsertBatch(ctx)
	if err != nil {
		pc.logger.Error("failed to insert batch", "error", err)
		pc.httpError(err, c)
		return
	}

	pc.logger.Info("insert batch")
	response.Created(c, "")
}

func (pc *CronController) httpError(err error, c *gin.Context) {
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		pc.logger.Warn("request deadline exceeded", "error", err)
		response.Custom(c, http.StatusGatewayTimeout, nil, "upstream service timed out")
	case errors.Is(err, context.Canceled):
		pc.logger.Warn("request was canceled", "error", err)
		response.Custom(c, http.StatusRequestTimeout, nil, "request was canceled by client")
	case errors.Is(err, service.ErrFailedToInsertBatchPrice):
		pc.logger.Warn("failed to insert batch", "error", err)
		response.Custom(c, http.StatusInternalServerError, nil, "failed to insert batch")
	case errors.Is(err, service.ErrFailedToGetPrice):
		pc.logger.Warn("failed to get price", "error", err)
		response.Custom(c, http.StatusInternalServerError, nil, "failed to get price")
	case errors.Is(err, currency.ErrCurrencyTooManyRequests):
		pc.logger.Warn("too many requests", "error", err)
		response.Custom(c, http.StatusTooManyRequests, nil, "too many requests")
	default:
		pc.logger.Error("internal server error", "error", err)
		response.InternalError(c)
	}
}
