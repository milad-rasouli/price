//go:build wireinject
// +build wireinject

package main

import (
	"github.com/milad-rasouli/price/internal/infrastructure/godotenv"
	"github.com/milad-rasouli/price/internal/infrastructure/postgresql"
	"github.com/milad-rasouli/price/internal/providers"
	"log/slog"

	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
	controller "github.com/milad-rasouli/price/internal/app/api/controllers"
	"github.com/milad-rasouli/price/internal/app/api/routes"
	"github.com/milad-rasouli/price/internal/repository/repository"
	"github.com/milad-rasouli/price/internal/service"
)

func wireApp(
	env *godotenv.Env,
	logger *slog.Logger,
	pg *postgresql.Postgres,
	pool *pgxpool.Pool,
) (*Boot, error) {
	panic(wire.Build(
		providers.ProviderSet,
		repository.ProviderSet,
		service.ProviderSet,
		controller.ProviderSet,
		routes.ProviderSet,
		wire.NewSet(NewBoot),
	))
}
