package service

import (
	"context"
	"errors"
	"testing"

	"github.com/AdhamHussam/wallet-ledge/internal/domain"
)

func TestCreateTransfer(t *testing.T) {
	ctx := context.Background()

	// Setup accounts
	acc1, err := testAccountService.CreateAccount(ctx, randomAccountRequest(500, domain.USD))
	if err != nil {
		t.Fatalf("failed to setup account 1: %v", err)
	}

	acc2, err := testAccountService.CreateAccount(ctx, randomAccountRequest(500, domain.USD))
	if err != nil {
		t.Fatalf("failed to setup account 2: %v", err)
	}

	accEUR, err := testAccountService.CreateAccount(ctx, randomAccountRequest(500, domain.EUR))
	if err != nil {
		t.Fatalf("failed to setup EUR account: %v", err)
	}

	// 1. Successful transfer
	res, err := testTransferService.CreateTransfer(ctx, domain.TransferRequest{
		FromAccountID: acc1.ID,
		ToAccountID:   acc2.ID,
		Amount:        150,
		Currency:      domain.USD,
	})
	if err != nil {
		t.Fatalf("expected successful transfer, got error: %v", err)
	}

	if res.Transfer.ID == 0 {
		t.Errorf("expected non-zero transfer ID")
	}
	if res.FromAccount.Balance != 350 {
		t.Errorf("expected sender balance 350, got %d", res.FromAccount.Balance)
	}
	if res.ToAccount.Balance != 650 {
		t.Errorf("expected recipient balance 650, got %d", res.ToAccount.Balance)
	}
	if res.FromEntry.Amount != -150 || res.ToEntry.Amount != 150 {
		t.Errorf("unexpected entry amounts: from=%d, to=%d", res.FromEntry.Amount, res.ToEntry.Amount)
	}

	// 2. Amount <= 0
	_, err = testTransferService.CreateTransfer(ctx, domain.TransferRequest{
		FromAccountID: acc1.ID,
		ToAccountID:   acc2.ID,
		Amount:        0,
		Currency:      domain.USD,
	})
	if !errors.Is(err, domain.ErrInvalidAmount) {
		t.Errorf("expected ErrInvalidAmount, got: %v", err)
	}

	// 3. Same account transfer
	_, err = testTransferService.CreateTransfer(ctx, domain.TransferRequest{
		FromAccountID: acc1.ID,
		ToAccountID:   acc1.ID,
		Amount:        50,
		Currency:      domain.USD,
	})
	if !errors.Is(err, domain.ErrSameAccountTransfer) {
		t.Errorf("expected ErrSameAccountTransfer, got: %v", err)
	}

	// 4. Unsupported currency
	_, err = testTransferService.CreateTransfer(ctx, domain.TransferRequest{
		FromAccountID: acc1.ID,
		ToAccountID:   acc2.ID,
		Amount:        50,
		Currency:      "BITCOIN",
	})
	if !errors.Is(err, domain.ErrInvalidCurrency) {
		t.Errorf("expected ErrInvalidCurrency, got: %v", err)
	}

	// 5. Currency mismatch (USD sender vs EUR receiver)
	_, err = testTransferService.CreateTransfer(ctx, domain.TransferRequest{
		FromAccountID: acc1.ID,
		ToAccountID:   accEUR.ID,
		Amount:        50,
		Currency:      domain.USD,
	})
	if !errors.Is(err, domain.ErrCurrencyMismatch) {
		t.Errorf("expected ErrCurrencyMismatch, got: %v", err)
	}

	// 6. Insufficient funds
	_, err = testTransferService.CreateTransfer(ctx, domain.TransferRequest{
		FromAccountID: acc1.ID,
		ToAccountID:   acc2.ID,
		Amount:        999999,
		Currency:      domain.USD,
	})
	if !errors.Is(err, domain.ErrInsufficientFunds) {
		t.Errorf("expected ErrInsufficientFunds, got: %v", err)
	}

	// 7. Non-existent sender
	_, err = testTransferService.CreateTransfer(ctx, domain.TransferRequest{
		FromAccountID: 999999999,
		ToAccountID:   acc2.ID,
		Amount:        10,
		Currency:      domain.USD,
	})
	if !errors.Is(err, domain.ErrAccountNotFound) {
		t.Errorf("expected ErrAccountNotFound, got: %v", err)
	}

	// 8. Non-existent receiver
	_, err = testTransferService.CreateTransfer(ctx, domain.TransferRequest{
		FromAccountID: acc1.ID,
		ToAccountID:   999999999,
		Amount:        10,
		Currency:      domain.USD,
	})
	if !errors.Is(err, domain.ErrAccountNotFound) {
		t.Errorf("expected ErrAccountNotFound, got: %v", err)
	}
}

func TestGetTransfer(t *testing.T) {
	ctx := context.Background()

	acc1, err := testAccountService.CreateAccount(ctx, randomAccountRequest(500, domain.GBP))
	if err != nil {
		t.Fatalf("failed to setup account 1: %v", err)
	}
	acc2, err := testAccountService.CreateAccount(ctx, randomAccountRequest(500, domain.GBP))
	if err != nil {
		t.Fatalf("failed to setup account 2: %v", err)
	}

	created, err := testTransferService.CreateTransfer(ctx, domain.TransferRequest{
		FromAccountID: acc1.ID,
		ToAccountID:   acc2.ID,
		Amount:        75,
		Currency:      domain.GBP,
	})
	if err != nil {
		t.Fatalf("failed to create transfer: %v", err)
	}

	// 1. Success
	fetched, err := testTransferService.GetTransfer(ctx, created.Transfer.ID)
	if err != nil {
		t.Fatalf("failed to get transfer: %v", err)
	}
	if fetched.ID != created.Transfer.ID || fetched.Amount != 75 {
		t.Errorf("fetched transfer mismatch: got %+v, want %+v", fetched, created.Transfer)
	}

	// 2. Not found
	_, err = testTransferService.GetTransfer(ctx, 999999999)
	if !errors.Is(err, domain.ErrTransferNotFound) {
		t.Errorf("expected ErrTransferNotFound, got: %v", err)
	}

	// 3. Invalid ID
	_, err = testTransferService.GetTransfer(ctx, -1)
	if !errors.Is(err, domain.ErrTransferNotFound) {
		t.Errorf("expected ErrTransferNotFound, got: %v", err)
	}
}
