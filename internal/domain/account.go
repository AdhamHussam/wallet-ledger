package domain

import "time"

// Account represents a financial ledger account.
type Account struct {
	ID        int64     `json:"id"`
	Owner     string    `json:"owner"`
	Balance   int64     `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}

// Entry records an immutable audit log entry for account balance changes.
type Entry struct {
	ID        int64     `json:"id"`
	AccountID int64     `json:"account_id"`
	Amount    int64     `json:"amount"` // Negative for debit, positive for credit
	CreatedAt time.Time `json:"created_at"`
}

// CreateAccountRequest contains data required to create a new account.
type CreateAccountRequest struct {
	Owner          string `json:"owner"`
	Currency       string `json:"currency"`
	InitialBalance int64  `json:"initial_balance,omitempty"`
}

// ListAccountsRequest contains pagination parameters for listing accounts.
type ListAccountsRequest struct {
	PageID   int32 `json:"page_id"`
	PageSize int32 `json:"page_size"`
}
