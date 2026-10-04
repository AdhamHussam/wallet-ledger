package db

import (
	"context"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/AdhamHussam/wallet-ledger/internal/domain"
)

// TestConcurrentMultiAccountTransfers verifies money conservation invariant
// across N accounts randomly sending money to each other concurrently.
// Invariant: Sum(Balance_initial) == Sum(Balance_final)
func TestConcurrentMultiAccountTransfers(t *testing.T) {
	ctx := context.Background()
	numAccounts := 6
	initialBalance := int64(1000)
	accounts := make([]Account, numAccounts)

	var totalInitialMoney int64
	for i := 0; i < numAccounts; i++ {
		accounts[i] = createRandomAccount(t, initialBalance)
		totalInitialMoney += initialBalance
	}

	numTransfers := 50
	transferAmount := int64(10)
	errChan := make(chan error, numTransfers)

	var wg sync.WaitGroup
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < numTransfers; i++ {
		fromIdx := r.Intn(numAccounts)
		toIdx := (fromIdx + 1 + r.Intn(numAccounts-1)) % numAccounts

		fromID := accounts[fromIdx].ID
		toID := accounts[toIdx].ID

		wg.Add(1)
		go func(from, to int64) {
			defer wg.Done()
			_, err := testStore.TransferTx(ctx, TransferTxParams{
				FromAccountID: from,
				ToAccountID:   to,
				Amount:        transferAmount,
			})
			errChan <- err
		}(fromID, toID)
	}

	wg.Wait()
	close(errChan)

	var successCount int
	for err := range errChan {
		if err != nil {
			// Some transfers might legitimately fail due to insufficient funds if an account is drawn down
			if err != domain.ErrInsufficientFunds {
				t.Fatalf("unexpected error during concurrent transfers: %v", err)
			}
		} else {
			successCount++
		}
	}

	t.Logf("Completed %d concurrent transfers (%d successful)", numTransfers, successCount)

	// Verify Conservation of Money: sum of all account balances must exactly equal initial sum
	var totalFinalMoney int64
	for _, acc := range accounts {
		updated, err := testStore.GetAccount(ctx, acc.ID)
		if err != nil {
			t.Fatalf("failed to retrieve account %d: %v", acc.ID, err)
		}
		if updated.Balance < 0 {
			t.Errorf("account %d has negative balance: %d", acc.ID, updated.Balance)
		}
		totalFinalMoney += updated.Balance
	}

	if totalFinalMoney != totalInitialMoney {
		t.Fatalf("Money conservation violated! Initial total: %d, Final total: %d (difference: %d)",
			totalInitialMoney, totalFinalMoney, totalInitialMoney-totalFinalMoney)
	}
}

// TestConcurrentOverdraftRace checks that when multiple concurrent goroutines try to withdraw
// more money than an account holds, only the valid number of transactions succeed
// and balance never goes negative.
func TestConcurrentOverdraftRace(t *testing.T) {
	ctx := context.Background()

	// Account has exactly 100
	sender := createRandomAccount(t, 100)
	receiver := createRandomAccount(t, 0)

	// 10 concurrent requests each trying to withdraw 50 (only 2 can ever succeed)
	concurrency := 10
	withdrawAmount := int64(50)

	var wg sync.WaitGroup
	errs := make(chan error, concurrency)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := testStore.TransferTx(ctx, TransferTxParams{
				FromAccountID: sender.ID,
				ToAccountID:   receiver.ID,
				Amount:        withdrawAmount,
			})
			errs <- err
		}()
	}

	wg.Wait()
	close(errs)

	var successes int
	var insufficientFunds int

	for err := range errs {
		if err == nil {
			successes++
		} else if err == domain.ErrInsufficientFunds {
			insufficientFunds++
		} else {
			t.Fatalf("unexpected transfer error: %v", err)
		}
	}

	if successes != 2 {
		t.Errorf("expected exactly 2 successful transfers, got %d", successes)
	}
	if insufficientFunds != concurrency-2 {
		t.Errorf("expected %d insufficient funds errors, got %d", concurrency-2, insufficientFunds)
	}

	updatedSender, err := testStore.GetAccount(ctx, sender.ID)
	if err != nil {
		t.Fatalf("failed to get sender: %v", err)
	}
	updatedReceiver, err := testStore.GetAccount(ctx, receiver.ID)
	if err != nil {
		t.Fatalf("failed to get receiver: %v", err)
	}

	if updatedSender.Balance != 0 {
		t.Errorf("expected sender balance 0, got %d", updatedSender.Balance)
	}
	if updatedReceiver.Balance != 100 {
		t.Errorf("expected receiver balance 100, got %d", updatedReceiver.Balance)
	}
}
