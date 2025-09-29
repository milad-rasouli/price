package controller

import (
	"context"
	"errors"
	"github.com/milad-rasouli/price/internal/app/api/dto"
	"github.com/milad-rasouli/price/internal/service"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/milad-rasouli/price/internal/app/api/response"
	"github.com/milad-rasouli/price/internal/repository/repository/price"
)

type PriceController struct {
	logger  *slog.Logger
	service service.PriceService
}

func NewPriceController(logger *slog.Logger, svc service.PriceService) *PriceController {
	return &PriceController{
		logger:  logger.With("layer", "PriceController"),
		service: svc,
	}
}

// GetHistory godoc
// @Summary Get historical cryptocurrency prices
// @Description Returns historical price data for a given symbol within a time range, optionally grouped by interval.
// @Tags prices
// @Accept json
// @Produce json
// @Param symbol query string true "Symbol (e.g., btc, eth)"
// @Param interval query string false "Interval (e.g., 1m, 5m, 1h, 1d)"
// @Param from query int false "Start time (unix timestamp)"
// @Param to query int false "End time (unix timestamp)"
// @Success 200 {object} response.Response[[]dto.HistoryRes]
// @Failure 400 {object} response.Response[any]
// @Failure 408 {object} response.Response[any]
// @Failure 504 {object} response.Response[any]
// @Failure 500 {object} response.Response[any]
// @Router /prices/history [get]
func (pc *PriceController) GetHistory(c *gin.Context) {
	req := &dto.HistoryReq{}
	if err := c.ShouldBindQuery(req); err != nil {
		response.BadRequest(c, "invalid query params: "+err.Error())
		return
	}

	// TODO: better validation using go-playground/validator

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	now := time.Now().Unix()
	if req.To == 0 {
		req.To = now
	}
	if req.From == 0 {
		req.From = req.To - 86400
	}

	history, err := pc.service.GetHistory(ctx, req)
	if err != nil {
		pc.logger.Error("failed to get history", "error", err, "symbol", req.Symbol)
		pc.httpError(err, c)
		return
	}

	pc.logger.Info("history fetched", "symbol", req.Symbol, "interval", req.Interval, "count", len(history))
	response.Ok(c, history, "")
}

// GetLatest godoc
// @Summary Get latest cryptocurrency price
// @Description Returns the latest stored price for a given symbol, including 24h change.
// @Tags prices
// @Accept json
// @Produce json
// @Param symbol query string true "Symbol (e.g., btc, eth)"
// @Success 200 {object} response.Response[dto.LatestRes]
// @Failure 400 {object} response.Response[any]
// @Failure 408 {object} response.Response[any]
// @Failure 504 {object} response.Response[any]
// @Failure 500 {object} response.Response[any]
// @Router /prices/latest [get]
func (pc *PriceController) GetLatest(c *gin.Context) {
	req := &dto.LatestReq{}
	if err := c.ShouldBindQuery(req); err != nil || req.Symbol == "" {
		response.BadRequest(c, "symbol is required")
		return
	}
	// TODO: better validation using go-playground/validator

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	latest, err := pc.service.GetLatest(ctx, req)
	if err != nil {
		pc.logger.Error("failed to get latest price", "error", err, "symbol", req.Symbol)
		pc.httpError(err, c)
		return
	}

	response.Ok(c, latest, "")
}

func (pc *PriceController) httpError(err error, c *gin.Context) {
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		pc.logger.Warn("request deadline exceeded", "error", err)
		response.Custom(c, http.StatusGatewayTimeout, nil, "upstream service timed out")
	case errors.Is(err, context.Canceled):
		pc.logger.Warn("request was canceled", "error", err)
		response.Custom(c, http.StatusRequestTimeout, nil, "request was canceled by client")
	case errors.Is(err, price.ErrPriceNotFound):
		response.NotFound(c)
	default:
		pc.logger.Error("internal server error", "error", err)
		response.InternalError(c)
	}
}
