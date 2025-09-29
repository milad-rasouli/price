package pgx

import (
	"context"
	"errors"
	"github.com/milad-rasouli/price/internal/app/api/dto"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/milad-rasouli/price/entity"
	"github.com/milad-rasouli/price/internal/repository/repository/price"
	"github.com/shopspring/decimal"
)

const (
	DefaultInterval = "1h"
	GetHistoryQuery = `
		SELECT EXTRACT(EPOCH FROM time_bucket($1, to_timestamp(time)))::BIGINT AS bucket,
			   symbol,
			   AVG(price) AS avg_price,
			   LAST(price, to_timestamp(time)) AS last_price
		FROM coin_prices
		WHERE symbol = $2
		  AND time BETWEEN $3 AND $4
		GROUP BY bucket, symbol
		ORDER BY bucket ASC
	`

	GetLatestQuery = `
		SELECT symbol, price, time
		FROM coin_prices
		WHERE symbol = $1
		ORDER BY time DESC
		LIMIT 1
	`

	GetBeforeTimeQuery = `
		SELECT price 
		FROM coin_prices 
		WHERE symbol = $1 AND time <= $2 
		ORDER BY time DESC LIMIT 1
	`
)

type PriceRepository struct {
	pool *pgxpool.Pool
}

func NewPriceRepository(pool *pgxpool.Pool) *PriceRepository {
	return &PriceRepository{pool: pool}
}

func (r *PriceRepository) BatchInsert(ctx context.Context, prices []*entity.Price) error {
	if len(prices) == 0 {
		return nil
	}

	rows := make([][]interface{}, len(prices))
	for i, p := range prices {
		rows[i] = []interface{}{p.Symbol, p.Price, p.Time}
	}

	_, err := r.pool.CopyFrom(
		ctx,
		pgx.Identifier{"coin_prices"},
		[]string{"symbol", "price", "time"},
		pgx.CopyFromRows(rows),
	)
	return err
}

func (r *PriceRepository) GetLatest(ctx context.Context, req *dto.LatestReq) (*dto.LatestRes, error) {
	var latest entity.Price
	err := r.pool.QueryRow(ctx, GetLatestQuery, req.Symbol).Scan(&latest.Symbol, &latest.Price, &latest.Time)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, price.ErrPriceNotFound
		}
		return nil, err
	}

	var price24h decimal.Decimal
	latesttime := time.Unix(latest.Time, 0)
	from := latesttime.Add(-24 * time.Hour).Unix()
	err = r.pool.QueryRow(ctx, GetBeforeTimeQuery, req.Symbol, from).Scan(&price24h)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	changePct := 0.0
	if !price24h.IsZero() {
		changePct = latest.Price.Sub(price24h).Div(price24h).Mul(decimal.NewFromInt(100)).InexactFloat64()
	}

	return &dto.LatestRes{
		Symbol:       latest.Symbol,
		Price:        latest.Price,
		Timestamp:    latest.Time,
		Change24HPct: changePct,
	}, nil
}

func (r *PriceRepository) GetHistory(ctx context.Context, req *dto.HistoryReq) ([]*dto.HistoryRes, error) {
	if req.Interval == "" {
		req.Interval = DefaultInterval
	}

	rows, err := r.pool.Query(ctx, GetHistoryQuery, req.Interval, req.Symbol, req.From, req.To)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*dto.HistoryRes
	for rows.Next() {
		var point dto.HistoryRes
		if err := rows.Scan(&point.StartedAt, &point.Symbol, &point.AvgPrice, &point.LastPrice); err != nil {
			return nil, err
		}
		result = append(result, &point)
	}

	if len(result) == 0 {
		return nil, price.ErrPriceNotFound
	}
	return result, nil
}
