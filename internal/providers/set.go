package providers

import (
	"github.com/google/wire"
	"github.com/milad-rasouli/price/internal/infrastructure/coingecko"
	"github.com/milad-rasouli/price/internal/providers/currency"
)

var ProviderSet = wire.NewSet(
	wire.Bind(new(currency.CurrencyProvider), new(*coingecko.CoinGecko)),
	coingecko.NewCoinGecko,
)
