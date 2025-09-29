package entity

import (
	"github.com/shopspring/decimal"
)

type Price struct {
	Symbol string          `json:"symbol"`
	Price  decimal.Decimal `json:"price"`
	Time   int64           `json:"time"`
}
