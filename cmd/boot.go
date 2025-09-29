package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/milad-rasouli/price/docs"
	"github.com/milad-rasouli/price/internal/app/api/routes"
	"github.com/milad-rasouli/price/internal/infrastructure/godotenv"
)

type Boot struct {
	env    *godotenv.Env
	logger *slog.Logger
	rts    []routes.Router
}

func NewBoot(
	env *godotenv.Env,
	logger *slog.Logger,
	rts ...routes.Router,
) *Boot {
	return &Boot{
		env:    env,
		logger: logger.With("layer", "boot"),
		rts:    rts,
	}
}

func (b *Boot) Boot() error {
	r := gin.Default()

	for _, router := range b.rts {
		router.SetupRoutes(r)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	addr := ":" + b.env.HTTPPort
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	serverErr := make(chan error, 1)
	go func() {
		b.logger.Info("starting HTTP server", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- fmt.Errorf("server failed: %w", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		return err
	case sig := <-quit:
		b.logger.Info("received shutdown signal", "signal", sig.String())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		b.logger.Error("server forced to shutdown", "error", err)
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	b.logger.Info("server exited cleanly")
	return nil
}
