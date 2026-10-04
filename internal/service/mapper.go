package service

import (
	"github.com/AdhamHussam/wallet-ledger/internal/domain"
	"github.com/AdhamHussam/wallet-ledger/internal/repository/db"
)

func toDomainAccount(a db.Account) domain.Account {
	return domain.Account{
		ID:        a.ID,
		Owner:     a.Owner,
		Balance:   a.Balance,
		Currency:  a.Currency,
		CreatedAt: a.CreatedAt.Time,
	}
}

func toDomainAccounts(accounts []db.Account) []domain.Account {
	result := make([]domain.Account, len(accounts))
	for i, a := range accounts {
		result[i] = toDomainAccount(a)
	}
	return result
}

func toDomainTransfer(t db.Transfer) domain.Transfer {
	return domain.Transfer{
		ID:            t.ID,
		FromAccountID: t.FromAccountID,
		ToAccountID:   t.ToAccountID,
		Amount:        t.Amount,
		CreatedAt:     t.CreatedAt.Time,
	}
}

func toDomainTransfers(transfers []db.Transfer) []domain.Transfer {
	result := make([]domain.Transfer, len(transfers))
	for i, t := range transfers {
		result[i] = toDomainTransfer(t)
	}
	return result
}

func toDomainEntry(e db.Entry) domain.Entry {
	return domain.Entry{
		ID:        e.ID,
		AccountID: e.AccountID,
		Amount:    e.Amount,
		CreatedAt: e.CreatedAt.Time,
	}
}
