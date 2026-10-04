package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/AdhamHussam/wallet-ledger/internal/domain"
)

func BenchmarkGetAccount(b *testing.B) {
	ctx := context.Background()
	acc, err := testStore.CreateAccount(ctx, CreateAccountParams{
		Owner:    "bench_user",
		Balance:  1000000,
		Currency: domain.USD,
	})
	if err != nil {
		b.Fatalf("failed to create account: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := testStore.GetAccount(ctx, acc.ID)
		if err != nil {
			b.Fatalf("failed to get account: %v", err)
		}
	}
}

func BenchmarkCreateAccount(b *testing.B) {
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := testStore.CreateAccount(ctx, CreateAccountParams{
			Owner:    fmt.Sprintf("bench_%d", i),
			Balance:  1000,
			Currency: domain.USD,
		})
		if err != nil {
			b.Fatalf("failed to create account: %v", err)
		}
	}
}

func BenchmarkTransferTx(b *testing.B) {
	ctx := context.Background()
	acc1, err := testStore.CreateAccount(ctx, CreateAccountParams{
		Owner:    "sender",
		Balance:  1000000000,
		Currency: domain.USD,
	})
	if err != nil {
		b.Fatalf("failed to create sender: %v", err)
	}

	acc2, err := testStore.CreateAccount(ctx, CreateAccountParams{
		Owner:    "receiver",
		Balance:  1000000000,
		Currency: domain.USD,
	})
	if err != nil {
		b.Fatalf("failed to create receiver: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := testStore.TransferTx(ctx, TransferTxParams{
			FromAccountID: acc1.ID,
			ToAccountID:   acc2.ID,
			Amount:        1,
		})
		if err != nil {
			b.Fatalf("failed to transfer: %v", err)
		}
	}
}
