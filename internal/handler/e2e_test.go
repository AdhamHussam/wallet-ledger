package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"

	"github.com/AdhamHussam/wallet-ledger/internal/domain"
	repository "github.com/AdhamHussam/wallet-ledger/internal/repository/db"
	"github.com/AdhamHussam/wallet-ledger/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	e2eServer  *httptest.Server
	e2eStore   repository.Store
	e2eAccount service.AccountService
)

func TestMain(m *testing.M) {
	dbURL := os.Getenv("DB_SOURCE")
	if dbURL == "" {
		dbURL = "postgres://postgres:secretpassword@localhost:5432/ledger_db?sslmode=disable"
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Printf("Cannot connect to DB for E2E tests: %v", err)
		os.Exit(0)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Printf("Cannot ping DB for E2E tests: %v", err)
		os.Exit(0)
	}

	e2eStore = repository.NewStore(pool)
	e2eAccount = service.NewAccountService(e2eStore)
	e2eTransfer := service.NewTransferService(e2eStore)

	router := NewRouter(e2eAccount, e2eTransfer)
	e2eServer = httptest.NewServer(router)
	defer e2eServer.Close()

	os.Exit(m.Run())
}

func TestE2EConcurrentTransfers(t *testing.T) {
	ctx := context.Background()

	// 1. Create two accounts via service
	acc1, err := e2eAccount.CreateAccount(ctx, domain.CreateAccountRequest{
		Owner:          "e2e_alice",
		Currency:       domain.USD,
		InitialBalance: 1000,
	})
	if err != nil {
		t.Fatalf("failed to create acc1: %v", err)
	}

	acc2, err := e2eAccount.CreateAccount(ctx, domain.CreateAccountRequest{
		Owner:          "e2e_bob",
		Currency:       domain.USD,
		InitialBalance: 1000,
	})
	if err != nil {
		t.Fatalf("failed to create acc2: %v", err)
	}

	// 2. Concurrently execute 20 HTTP POST /transfers (10 alice->bob, 10 bob->alice)
	concurrency := 20
	var wg sync.WaitGroup
	statusCodes := make(chan int, concurrency)

	client := e2eServer.Client()

	for i := 0; i < concurrency; i++ {
		fromID := acc1.ID
		toID := acc2.ID
		if i%2 == 1 {
			fromID = acc2.ID
			toID = acc1.ID
		}

		wg.Add(1)
		go func(from, to int64) {
			defer wg.Done()
			payload, _ := json.Marshal(domain.TransferRequest{
				FromAccountID: from,
				ToAccountID:   to,
				Amount:        25,
				Currency:      domain.USD,
			})

			resp, err := client.Post(
				fmt.Sprintf("%s/transfers", e2eServer.URL),
				"application/json",
				bytes.NewReader(payload),
			)
			if err != nil {
				t.Errorf("HTTP post error: %v", err)
				return
			}
			defer resp.Body.Close()
			_, _ = io.ReadAll(resp.Body)
			statusCodes <- resp.StatusCode
		}(fromID, toID)
	}

	wg.Wait()
	close(statusCodes)

	for code := range statusCodes {
		if code != http.StatusCreated {
			t.Errorf("expected status 201 Created, got %d", code)
		}
	}

	// 3. Verify final account balances equal original balance (since 10 sent and 10 received)
	final1, err := e2eAccount.GetAccount(ctx, acc1.ID)
	if err != nil {
		t.Fatalf("failed to get account 1: %v", err)
	}
	final2, err := e2eAccount.GetAccount(ctx, acc2.ID)
	if err != nil {
		t.Fatalf("failed to get account 2: %v", err)
	}

	if final1.Balance != 1000 {
		t.Errorf("expected final balance 1000 for alice, got %d", final1.Balance)
	}
	if final2.Balance != 1000 {
		t.Errorf("expected final balance 1000 for bob, got %d", final2.Balance)
	}
}
