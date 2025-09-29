package price

import (
	"context"
	"errors"
	"github.com/milad-rasouli/price/entity"
	"github.com/milad-rasouli/price/internal/app/api/dto"
)

var (
	ErrPriceNotFound = errors.New("price not found")
)

//go:generate mockgen -source=price.go -destination=../../../../mock/repository/price/price.go
type PriceRepository interface {
	BatchInsert(ctx context.Context, p []*entity.Price) error
	GetLatest(ctx context.Context, req *dto.LatestReq) (*dto.LatestRes, error)
	GetHistory(ctx context.Context, req *dto.HistoryReq) ([]*dto.HistoryRes, error)
}
