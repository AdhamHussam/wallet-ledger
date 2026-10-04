package domain

import "errors"

var (
	// ErrAccountNotFound is returned when an account cannot be found.
	ErrAccountNotFound = errors.New("account not found")

	// ErrInsufficientFunds is returned when an account does not have enough balance for a debit/transfer.
	ErrInsufficientFunds = errors.New("insufficient funds")

	// ErrInvalidAmount is returned when an amount is negative or zero.
	ErrInvalidAmount = errors.New("amount must be greater than zero")

	// ErrSameAccountTransfer is returned when source and destination accounts are the same.
	ErrSameAccountTransfer = errors.New("cannot transfer money to the same account")

	// ErrCurrencyMismatch is returned when currencies between accounts or transfer requests do not match.
	ErrCurrencyMismatch = errors.New("currency mismatch between accounts")

	// ErrInvalidCurrency is returned when a currency code is not supported.
	ErrInvalidCurrency = errors.New("invalid or unsupported currency")
)
