package main

import (
	"context"
	"flag"
	"github.com/milad-rasouli/price/cmd/cron"
	"log/slog"
	"os"
	"time"

	"github.com/milad-rasouli/price/internal/infrastructure/godotenv"
	"github.com/milad-rasouli/price/internal/infrastructure/postgresql"
)

func main() {
	env := godotenv.NewEnv()
	logger := initSlogLogger(env)
	cronflag := flag.Bool("cron", false, "Trigger cron jobs")
	flag.Parse()
	if *cronflag {
		if err := cron.Run(env, logger); err != nil {
			logger.Error("failed to run cron", "error", err)
			os.Exit(1)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pg := postgresql.NewPostgre(env)
	err := pg.Setup(ctx)
	if err != nil {
		logger.Error("failed to initialize PostgreSQL Primary", "error", err)
		os.Exit(1)
	}
	defer pg.Close()

	err = pg.HealthCheck(ctx)
	if err != nil {
		logger.Error("failed to check health", "error", err)
		os.Exit(1)
	}
	pool := pg.Pool

	boot, err := wireApp(env, logger, pg, pool)
	if err != nil {
		logger.Error("failed to setup app", "error", err)
		os.Exit(1)
	}
	err = boot.Boot()
	if err != nil {
		logger.Error("failed to start app", "error", err)
		os.Exit(1)
	}

}

func initSlogLogger(e *godotenv.Env) *slog.Logger {
	logLevel := slog.LevelDebug
	if e.Environment == "production" {
		logLevel = slog.LevelWarn
	}
	//development
	slogHandlerOptions := &slog.HandlerOptions{
		AddSource: true,
		Level:     logLevel,
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, slogHandlerOptions))
	slog.SetDefault(logger)

	return logger
}
