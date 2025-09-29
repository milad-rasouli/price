package cron

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/milad-rasouli/price/internal/infrastructure/godotenv"
)

func Run(env *godotenv.Env, logger *slog.Logger) error {
	lg := logger.With("method", "cron.Run")

	interval := time.Duration(env.ReadCoinInterval) * time.Second

	url := fmt.Sprintf("http://localhost:%s/cron/update-prices", env.HTTPPort)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	lg.Info("cron started", "interval", interval, "url", url)

	for {
		select {
		case <-ticker.C:
			callUpdatePrices(lg, url)

		case <-ctx.Done():
			lg.Info("cron shutting down gracefully...")
			return nil
		}
	}
}

func callUpdatePrices(lg *slog.Logger, url string) {
	resp, err := http.Get(url)
	if err != nil {
		lg.Error("failed to call update-prices", "error", err)
		return
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			lg.Error("failed to close response body", "error", cerr)
		}
	}()

	if resp.StatusCode == http.StatusOK {
		lg.Info("successfully updated prices", "status", resp.Status)
	} else {
		lg.Error("unexpected response from update-prices",
			"status", resp.Status,
			"code", resp.StatusCode,
		)
	}
}
