package domain

import "time"

// Transfer represents a financial transfer between two accounts.
type Transfer struct {
	ID            int64     `json:"id"`
	FromAccountID int64     `json:"from_account_id"`
	ToAccountID   int64     `json:"to_account_id"`
	Amount        int64     `json:"amount"`
	CreatedAt     time.Time `json:"created_at"`
}

// TransferRequest contains data required to initiate a money transfer.
type TransferRequest struct {
	FromAccountID int64  `json:"from_account_id"`
	ToAccountID   int64  `json:"to_account_id"`
	Amount        int64  `json:"amount"`
	Currency      string `json:"currency"`
}

// TransferResponse details the outcome of a money transfer transaction.
type TransferResponse struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// ListTransfersRequest contains pagination parameters for listing transfers.
type ListTransfersRequest struct {
	PageID   int32 `json:"page_id"`
	PageSize int32 `json:"page_size"`
}
