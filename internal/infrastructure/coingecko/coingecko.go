package coingecko

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/milad-rasouli/price/entity"
	"github.com/milad-rasouli/price/internal/providers/currency"
	"github.com/shopspring/decimal"
	"log/slog"
	"net/http"
	"time"
)

type CoinGecko struct {
	client  *http.Client
	baseURL string
	logger  *slog.Logger
}

func NewCoinGecko(logger *slog.Logger) *CoinGecko {
	return &CoinGecko{
		client:  &http.Client{Timeout: 5 * time.Second},
		baseURL: "https://api.coingecko.com/api/v3",
		logger:  logger.With("provider", "coingecko"),
	}
}

type coinResponse struct {
	Symbol       string  `json:"symbol"`
	CurrentPrice float64 `json:"current_price"`
	LastUpdated  string  `json:"last_updated"`
}

func (c *CoinGecko) Get(ctx context.Context, page, limit uint32) ([]*entity.Price, error) {
	url := fmt.Sprintf(
		"%s/coins/markets?vs_currency=usd&order=market_cap_desc&per_page=%d&page=%d&sparkline=false&price_change_percentage=24h",
		c.baseURL, limit, page,
	)

	c.logger.Info("fetching prices from coingecko", "url", url, "page", page, "limit", limit)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		c.logger.Error("failed to create request", "error", err)
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		c.logger.Error("failed to call coingecko API", "error", err)
		return nil, err
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			c.logger.Warn("failed to close response body", "error", cerr)
		}
	}()

	switch resp.StatusCode {
	case http.StatusOK:
		// continue
	case http.StatusTooManyRequests:
		c.logger.Error("rate limit exceeded from coingecko", "status", resp.StatusCode)
		return nil, currency.ErrCurrencyTooManyRequests
	default:
		c.logger.Error("unexpected status code", "status", resp.StatusCode)
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var coins []coinResponse
	if err := json.NewDecoder(resp.Body).Decode(&coins); err != nil {
		c.logger.Error("failed to decode coingecko response", "error", err)
		return nil, err
	}

	if len(coins) == 0 {
		c.logger.Warn("no currencies returned from coingecko", "page", page, "limit", limit)
		return nil, currency.ErrCurrencyNotFound
	}

	result := make([]*entity.Price, 0, len(coins))
	for _, coin := range coins {
		t, err := time.Parse(time.RFC3339, coin.LastUpdated)
		unixTime := time.Now().Unix()
		if err == nil {
			unixTime = t.Unix()
		} else {
			c.logger.Warn("invalid last_updated format, using now()", "symbol", coin.Symbol, "last_updated", coin.LastUpdated)
		}

		result = append(result, &entity.Price{
			Symbol: coin.Symbol,
			Price:  decimal.NewFromFloat(coin.CurrentPrice),
			Time:   unixTime,
		})
	}

	c.logger.Info("fetched prices successfully", "count", len(result))
	return result, nil
}
