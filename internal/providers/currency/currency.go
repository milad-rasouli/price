package currency

import (
	"context"
	"errors"
	"github.com/milad-rasouli/price/entity"
)

var (
	ErrCurrencyNotFound        = errors.New("currency not found")
	ErrCurrencyTooManyRequests = errors.New("too many requests")
)

//go:generate mockgen -source=currency.go -destination=../../../mock/providers/currency/currency.go
type CurrencyProvider interface {
	Get(ctx context.Context, page, limit uint32) ([]*entity.Price, error)
}
