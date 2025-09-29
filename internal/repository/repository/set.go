package repository

import (
	"github.com/google/wire"
	"github.com/milad-rasouli/price/internal/repository/repository/price"
	"github.com/milad-rasouli/price/internal/repository/repository/price/pgx"
)

var ProviderSet = wire.NewSet(
	wire.Bind(new(price.PriceRepository), new(*pgx.PriceRepository)),
	pgx.NewPriceRepository,
)
