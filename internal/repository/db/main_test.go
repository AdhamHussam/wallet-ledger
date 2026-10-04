package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testStore Store

func TestMain(m *testing.M) {
	dbURL := os.Getenv("DB_SOURCE")
	if dbURL == "" {
		dbURL = "postgres://postgres:secretpassword@localhost:5432/ledger_db?sslmode=disable"
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Printf("Unable to connect to database for integration tests: %v", err)
		os.Exit(0)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Printf("Database ping failed: %v", err)
		os.Exit(0)
	}

	testStore = NewStore(pool)
	os.Exit(m.Run())
}
