package domain

import "time"

type Account struct {
	ID        int64     `json:"id"`
	Owner     string    `json:"owner"`
	Balance   int64     `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}

type Entry struct {
	ID        int64     `json:"id"`
	AccountID int64     `json:"account_id"`
	Amount    int64     `json:"amount"` // Negative for debit, positive for credit
	CreatedAt time.Time `json:"created_at"`
}

type CreateAccountRequest struct {
	Owner          string `json:"owner"`
	Currency       string `json:"currency"`
	InitialBalance int64  `json:"initial_balance,omitempty"`
}

type ListAccountsRequest struct {
	PageID   int32 `json:"page_id"`
	PageSize int32 `json:"page_size"`
}
