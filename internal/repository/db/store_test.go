package db

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/AdhamHussam/wallet-ledger/internal/domain"
)

func createRandomAccount(t *testing.T, balance int64) Account {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	arg := CreateAccountParams{
		Owner:    fmt.Sprintf("user_%d", r.Intn(1000000)),
		Balance:  balance,
		Currency: domain.USD,
	}

	account, err := testStore.CreateAccount(context.Background(), arg)
	if err != nil {
		t.Fatalf("failed to create random account: %v", err)
	}

	return account
}

func TestTransferTx(t *testing.T) {
	account1 := createRandomAccount(t, 1000)
	account2 := createRandomAccount(t, 1000)

	n := 5
	amount := int64(10)

	errs := make(chan error, n)
	results := make(chan TransferTxResult, n)

	for i := 0; i < n; i++ {
		go func() {
			result, err := testStore.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		if err != nil {
			t.Fatalf("unexpected transfer error: %v", err)
		}

		result := <-results
		if result.Transfer.ID == 0 {
			t.Errorf("expected transfer ID to be non-zero")
		}
		if result.Transfer.FromAccountID != account1.ID {
			t.Errorf("expected from_account_id %d, got %d", account1.ID, result.Transfer.FromAccountID)
		}
		if result.Transfer.ToAccountID != account2.ID {
			t.Errorf("expected to_account_id %d, got %d", account2.ID, result.Transfer.ToAccountID)
		}
		if result.Transfer.Amount != amount {
			t.Errorf("expected amount %d, got %d", amount, result.Transfer.Amount)
		}

		// Check entries
		if result.FromEntry.AccountID != account1.ID || result.FromEntry.Amount != -amount {
			t.Errorf("invalid from_entry: %+v", result.FromEntry)
		}
		if result.ToEntry.AccountID != account2.ID || result.ToEntry.Amount != amount {
			t.Errorf("invalid to_entry: %+v", result.ToEntry)
		}

		// Check account balances
		diff1 := account1.Balance - result.FromAccount.Balance
		diff2 := result.ToAccount.Balance - account2.Balance
		if diff1 != diff2 {
			t.Errorf("balance diff mismatch: diff1=%d, diff2=%d", diff1, diff2)
		}
		if diff1 <= 0 || diff1%amount != 0 {
			t.Errorf("invalid balance difference: %d", diff1)
		}

		k := int(diff1 / amount)
		if k < 1 || k > n {
			t.Errorf("k out of bounds: %d", k)
		}
		if existed[k] {
			t.Errorf("duplicate step detected: %d", k)
		}
		existed[k] = true
	}

	// Verify final balances in database
	updatedAccount1, err := testStore.GetAccount(context.Background(), account1.ID)
	if err != nil {
		t.Fatalf("failed to get account 1: %v", err)
	}
	updatedAccount2, err := testStore.GetAccount(context.Background(), account2.ID)
	if err != nil {
		t.Fatalf("failed to get account 2: %v", err)
	}

	expectedBalance1 := account1.Balance - int64(n)*amount
	expectedBalance2 := account2.Balance + int64(n)*amount

	if updatedAccount1.Balance != expectedBalance1 {
		t.Errorf("expected account 1 balance %d, got %d", expectedBalance1, updatedAccount1.Balance)
	}
	if updatedAccount2.Balance != expectedBalance2 {
		t.Errorf("expected account 2 balance %d, got %d", expectedBalance2, updatedAccount2.Balance)
	}
}

func TestTransferTxDeadlock(t *testing.T) {
	account1 := createRandomAccount(t, 1000)
	account2 := createRandomAccount(t, 1000)

	n := 10
	amount := int64(10)
	errs := make(chan error, n)

	// Run n transfers: half from 1 -> 2, half from 2 -> 1 concurrently
	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		go func(from, to int64) {
			_, err := testStore.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: from,
				ToAccountID:   to,
				Amount:        amount,
			})
			errs <- err
		}(fromAccountID, toAccountID)
	}

	for i := 0; i < n; i++ {
		err := <-errs
		if err != nil {
			t.Fatalf("transfer transaction failed (possible deadlock): %v", err)
		}
	}

	// Verify final balances remain equal since n/2 transfers went each way
	updatedAccount1, err := testStore.GetAccount(context.Background(), account1.ID)
	if err != nil {
		t.Fatalf("failed to get account 1: %v", err)
	}
	updatedAccount2, err := testStore.GetAccount(context.Background(), account2.ID)
	if err != nil {
		t.Fatalf("failed to get account 2: %v", err)
	}

	if updatedAccount1.Balance != account1.Balance {
		t.Errorf("expected account 1 balance %d, got %d", account1.Balance, updatedAccount1.Balance)
	}
	if updatedAccount2.Balance != account2.Balance {
		t.Errorf("expected account 2 balance %d, got %d", account2.Balance, updatedAccount2.Balance)
	}
}

func TestTransferTxInsufficientFunds(t *testing.T) {
	account1 := createRandomAccount(t, 50)
	account2 := createRandomAccount(t, 100)

	// Attempt transfer greater than balance
	_, err := testStore.TransferTx(context.Background(), TransferTxParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        100,
	})

	if !errors.Is(err, domain.ErrInsufficientFunds) {
		t.Fatalf("expected ErrInsufficientFunds, got: %v", err)
	}

	// Verify balances remained untouched
	check1, err := testStore.GetAccount(context.Background(), account1.ID)
	if err != nil {
		t.Fatalf("failed to get account 1: %v", err)
	}
	check2, err := testStore.GetAccount(context.Background(), account2.ID)
	if err != nil {
		t.Fatalf("failed to get account 2: %v", err)
	}

	if check1.Balance != 50 {
		t.Errorf("expected balance to stay 50, got %d", check1.Balance)
	}
	if check2.Balance != 100 {
		t.Errorf("expected balance to stay 100, got %d", check2.Balance)
	}
}
