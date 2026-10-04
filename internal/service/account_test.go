package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/AdhamHussam/wallet-ledge/internal/domain"
)

func randomAccountRequest(balance int64, currency string) domain.CreateAccountRequest {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return domain.CreateAccountRequest{
		Owner:          fmt.Sprintf("user_%d", r.Intn(1000000)),
		Currency:       currency,
		InitialBalance: balance,
	}
}

func TestCreateAccount(t *testing.T) {
	ctx := context.Background()

	// 1. Success case
	req := randomAccountRequest(500, domain.USD)
	account, err := testAccountService.CreateAccount(ctx, req)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if account.ID == 0 {
		t.Errorf("expected non-zero account ID")
	}
	if account.Owner != req.Owner {
		t.Errorf("expected owner %s, got %s", req.Owner, account.Owner)
	}
	if account.Balance != 500 {
		t.Errorf("expected balance 500, got %d", account.Balance)
	}
	if account.Currency != domain.USD {
		t.Errorf("expected currency USD, got %s", account.Currency)
	}

	// 2. Empty owner error
	_, err = testAccountService.CreateAccount(ctx, domain.CreateAccountRequest{
		Owner:    "   ",
		Currency: domain.USD,
	})
	if !errors.Is(err, domain.ErrInvalidOwner) {
		t.Errorf("expected ErrInvalidOwner, got: %v", err)
	}

	// 3. Invalid currency error
	_, err = testAccountService.CreateAccount(ctx, domain.CreateAccountRequest{
		Owner:    "valid_owner",
		Currency: "INVALID",
	})
	if !errors.Is(err, domain.ErrInvalidCurrency) {
		t.Errorf("expected ErrInvalidCurrency, got: %v", err)
	}

	// 4. Negative initial balance error
	_, err = testAccountService.CreateAccount(ctx, domain.CreateAccountRequest{
		Owner:          "valid_owner",
		Currency:       domain.USD,
		InitialBalance: -10,
	})
	if !errors.Is(err, domain.ErrInvalidAmount) {
		t.Errorf("expected ErrInvalidAmount, got: %v", err)
	}
}

func TestGetAccount(t *testing.T) {
	ctx := context.Background()

	// Create an account first
	created, err := testAccountService.CreateAccount(ctx, randomAccountRequest(200, domain.EUR))
	if err != nil {
		t.Fatalf("failed to create account: %v", err)
	}

	// 1. Success
	fetched, err := testAccountService.GetAccount(ctx, created.ID)
	if err != nil {
		t.Fatalf("failed to get account: %v", err)
	}
	if fetched.ID != created.ID || fetched.Balance != created.Balance || fetched.Currency != created.Currency {
		t.Errorf("fetched account does not match: got %+v, want %+v", fetched, created)
	}

	// 2. Not found
	_, err = testAccountService.GetAccount(ctx, 999999999)
	if !errors.Is(err, domain.ErrAccountNotFound) {
		t.Errorf("expected ErrAccountNotFound, got: %v", err)
	}

	// 3. Invalid ID
	_, err = testAccountService.GetAccount(ctx, -1)
	if !errors.Is(err, domain.ErrAccountNotFound) {
		t.Errorf("expected ErrAccountNotFound, got: %v", err)
	}
}

func TestListAccounts(t *testing.T) {
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		_, err := testAccountService.CreateAccount(ctx, randomAccountRequest(100, domain.USD))
		if err != nil {
			t.Fatalf("failed to create account: %v", err)
		}
	}

	accounts, err := testAccountService.ListAccounts(ctx, domain.ListAccountsRequest{
		PageID:   1,
		PageSize: 5,
	})
	if err != nil {
		t.Fatalf("failed to list accounts: %v", err)
	}

	if len(accounts) < 5 {
		t.Errorf("expected at least 5 accounts, got %d", len(accounts))
	}
}
