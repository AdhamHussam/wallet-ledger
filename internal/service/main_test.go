package service

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/AdhamHussam/wallet-ledge/internal/repository/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	testStore           db.Store
	testAccountService  AccountService
	testTransferService TransferService
)

func TestMain(m *testing.M) {
	dbURL := os.Getenv("DB_SOURCE")
	if dbURL == "" {
		dbURL = "postgres://postgres:secretpassword@localhost:5432/ledger_db?sslmode=disable"
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Printf("Unable to connect to database for service tests: %v", err)
		os.Exit(0)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Printf("Database ping failed: %v", err)
		os.Exit(0)
	}

	testStore = db.NewStore(pool)
	testAccountService = NewAccountService(testStore)
	testTransferService = NewTransferService(testStore)

	os.Exit(m.Run())
}
