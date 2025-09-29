package postgresql

import (
	"context"
	"fmt"
	"github.com/milad-rasouli/price/internal/infrastructure/godotenv"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Postgres struct {
	env  *godotenv.Env
	Pool *pgxpool.Pool
}

func NewPostgre(env *godotenv.Env) *Postgres {
	return &Postgres{
		env: env,
	}
}

func (p *Postgres) Setup(ctx context.Context) error {
	pool, err := pgxpool.New(ctx, p.env.DatabaseHost)
	if err != nil {
		return err
	}
	p.Pool = pool
	return nil
}

func (p *Postgres) HealthCheck(ctx context.Context) error {
	if p.Pool == nil {
		return fmt.Errorf("connection Pool not initialized")
	}

	err := p.Pool.Ping(ctx)
	if err != nil {
		return fmt.Errorf("failed to ping connection Pool: %w", err)
	}

	return nil
}

func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
