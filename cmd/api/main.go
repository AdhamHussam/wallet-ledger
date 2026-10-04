package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AdhamHussam/wallet-ledger/db"
	"github.com/AdhamHussam/wallet-ledger/internal/config"
	"github.com/AdhamHussam/wallet-ledger/internal/handler"
	repository "github.com/AdhamHussam/wallet-ledger/internal/repository/db"
	"github.com/AdhamHussam/wallet-ledger/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	if cfg.RunMigrationOnBoot {
		slog.Info("running database migrations...")
		if err := db.RunMigrations(cfg.DBSource); err != nil {
			log.Fatalf("Failed to apply database migrations: %v", err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	poolCfg, err := pgxpool.ParseConfig(cfg.DBSource)
	if err != nil {
		log.Fatalf("Unable to parse DB connection string: %v", err)
	}
	poolCfg.MaxConns = 25
	poolCfg.MinConns = 5
	poolCfg.MaxConnLifetime = 1 * time.Hour
	poolCfg.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Cannot reach database: %v", err)
	}
	slog.Info("connected to PostgreSQL", "environment", cfg.Environment)

	store := repository.NewStore(pool)
	accountSvc := service.NewAccountService(store)
	transferSvc := service.NewTransferService(store)
	router := handler.NewRouter(accountSvc, transferSvc)

	server := &http.Server{
		Addr:         cfg.HTTPServerAddress,
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	serverErrors := make(chan error, 1)
	go func() {
		slog.Info("server starting", "addr", cfg.HTTPServerAddress)
		serverErrors <- server.ListenAndServe()
	}()

	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %v", err)
		}
	case sig := <-shutdownSignal:
		slog.Info("shutdown signal received", "signal", sig.String())

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			slog.Error("graceful shutdown failed, forcing close", "error", err)
			_ = server.Close()
		} else {
			slog.Info("server gracefully stopped")
		}
	}
}
