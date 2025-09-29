package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/milad-rasouli/price/entity"
	"github.com/milad-rasouli/price/internal/app/api/dto"
	"github.com/milad-rasouli/price/internal/providers/currency"
	"github.com/milad-rasouli/price/internal/repository/repository/price"
	"log/slog"
	"time"
)

var (
	ErrFailedToGetPrice         = errors.New("failed to get price")
	ErrFailedToInsertBatchPrice = errors.New("failed to insert batch price")
)

const (
	MaxRetry     int           = 3
	BackoffDelay time.Duration = 1 * time.Second
)

//go:generate mockgen -source=price.go -destination=../../mock/service/price/price.go
type PriceService interface {
	InsertBatch(ctx context.Context) error
	GetLatest(ctx context.Context, req *dto.LatestReq) (*dto.LatestRes, error)
	GetHistory(ctx context.Context, req *dto.HistoryReq) ([]*dto.HistoryRes, error)
}

type priceService struct {
	logger           *slog.Logger
	repo             price.PriceRepository
	currencyProvider currency.CurrencyProvider
}

func NewPriceService(
	logger *slog.Logger,
	repo price.PriceRepository,
	currencyProvider currency.CurrencyProvider,
) PriceService {
	return &priceService{
		logger:           logger.With("Layer", "PriceService"),
		repo:             repo,
		currencyProvider: currencyProvider,
	}
}

func (s *priceService) InsertBatch(ctx context.Context) error {
	var (
		lg     = s.logger.With("method", "InsertBatch")
		prices []*entity.Price
		err    error
	)
	for attempt := 1; attempt <= MaxRetry; attempt++ {
		prices, err = s.currencyProvider.Get(ctx, 1, 3)
		if err == nil {
			break
		}
		if errors.Is(err, currency.ErrCurrencyTooManyRequests) {
			lg.Warn("Get currency is too many requests", "error", err)
			return fmt.Errorf("%w (attempt %d)", err, attempt)
		}
		if attempt == MaxRetry {
			lg.Error("failed to get prices", "attempt", attempt, "error", err)
			return fmt.Errorf("%w: %s", ErrFailedToGetPrice, err)
		}

		lg.Warn("failed to get prices", "attempt", attempt, "error", err)
		select {
		case <-time.After(BackoffDelay):
		case <-ctx.Done():
			return ctx.Err()
		}

	}
	if err := s.repo.BatchInsert(ctx, prices); err != nil {
		lg.Error("failed to batch insert prices", "error", err)
		return fmt.Errorf("%w: %s", ErrFailedToInsertBatchPrice, err)
	}
	lg.Info("successfully inserted batch of prices", "count", len(prices))
	return nil
}

func (s *priceService) GetLatest(ctx context.Context, req *dto.LatestReq) (*dto.LatestRes, error) {
	lg := s.logger.With("method", "GetLatest")
	latest, err := s.repo.GetLatest(ctx, req)
	if err != nil {
		lg.Error("failed to get latest price", "symbol", req.Symbol, "error", err)
		return nil, err
	}

	lg.Info("fetched latest price", "symbol", req.Symbol, "price", latest.Price, "change_24h_pct", latest.Change24HPct)
	return latest, nil
}

func (s *priceService) GetHistory(ctx context.Context, req *dto.HistoryReq) ([]*dto.HistoryRes, error) {
	lg := s.logger.With("method", "GetHistory")

	history, err := s.repo.GetHistory(ctx, req)
	if err != nil {
		lg.Error("failed to fetch price history", "symbol", req.Symbol, "error", err)
		return nil, err
	}

	lg.Info("fetched price history", "symbol", req.Symbol, "points", len(history))
	return history, nil
}
