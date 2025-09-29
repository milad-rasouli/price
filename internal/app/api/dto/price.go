package dto

import (
	"github.com/shopspring/decimal"
)

type HistoryReq struct {
	Symbol   string `form:"symbol" binding:"required"`
	Interval string `form:"interval"` // e.g. "1m", "5m", "1h", "1d"
	From     int64  `form:"from"`
	To       int64  `form:"to"`
}

type LatestReq struct {
	Symbol string `form:"symbol" binding:"required"`
}

type LatestRes struct {
	Symbol       string          `json:"symbol"`
	Price        decimal.Decimal `json:"price"`
	Timestamp    int64           `json:"timestamp"`
	Change24HPct float64         `json:"change_24h_pct"`
}

type HistoryRes struct {
	StartedAt int64           `json:"startedAt"`
	Symbol    string          `json:"symbol"`
	AvgPrice  decimal.Decimal `json:"avg_price"`
	LastPrice decimal.Decimal `json:"last_price"`
}
