package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/AdhamHussam/wallet-ledge/internal/handler"
	"github.com/AdhamHussam/wallet-ledge/internal/repository/db"
	"github.com/AdhamHussam/wallet-ledge/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dbURL := "postgres://postgres:secretpassword@localhost:5432/ledger_db?sslmode=disable"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Cannot reach database: %v\n", err)
	}
	slog.Info("connected to PostgreSQL")

	store := db.NewStore(pool)
	accountSvc := service.NewAccountService(store)
	transferSvc := service.NewTransferService(store)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      handler.NewRouter(accountSvc, transferSvc),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	slog.Info("server running", "addr", "http://localhost:8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v\n", err)
	}
}
