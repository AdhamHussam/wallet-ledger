package service

import (
	"context"
	"errors"
	"strings"

	"github.com/AdhamHussam/wallet-ledger/internal/domain"
	"github.com/AdhamHussam/wallet-ledger/internal/repository/db"
	"github.com/jackc/pgx/v5"
)

// TransferService defines business operations on money transfers.
type TransferService interface {
	CreateTransfer(ctx context.Context, req domain.TransferRequest) (domain.TransferResponse, error)
	GetTransfer(ctx context.Context, id int64) (domain.Transfer, error)
	ListTransfers(ctx context.Context, req domain.ListTransfersRequest) ([]domain.Transfer, error)
}

type transferService struct {
	store db.Store
}

// NewTransferService creates a new TransferService.
func NewTransferService(store db.Store) TransferService {
	return &transferService{store: store}
}

func (s *transferService) CreateTransfer(ctx context.Context, req domain.TransferRequest) (domain.TransferResponse, error) {
	if req.Amount <= 0 {
		return domain.TransferResponse{}, domain.ErrInvalidAmount
	}

	if req.FromAccountID == req.ToAccountID {
		return domain.TransferResponse{}, domain.ErrSameAccountTransfer
	}

	if req.FromAccountID <= 0 || req.ToAccountID <= 0 {
		return domain.TransferResponse{}, domain.ErrAccountNotFound
	}

	currency := strings.ToUpper(strings.TrimSpace(req.Currency))
	if !domain.IsSupportedCurrency(currency) {
		return domain.TransferResponse{}, domain.ErrInvalidCurrency
	}

	// 1. Verify source account exists and check currency & balance
	fromAccount, err := s.store.GetAccount(ctx, req.FromAccountID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.TransferResponse{}, domain.ErrAccountNotFound
		}
		return domain.TransferResponse{}, err
	}

	if fromAccount.Currency != currency {
		return domain.TransferResponse{}, domain.ErrCurrencyMismatch
	}

	if fromAccount.Balance < req.Amount {
		return domain.TransferResponse{}, domain.ErrInsufficientFunds
	}

	// 2. Verify destination account exists and check currency
	toAccount, err := s.store.GetAccount(ctx, req.ToAccountID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.TransferResponse{}, domain.ErrAccountNotFound
		}
		return domain.TransferResponse{}, err
	}

	if toAccount.Currency != currency {
		return domain.TransferResponse{}, domain.ErrCurrencyMismatch
	}

	// 3. Execute atomic transaction in the store
	result, err := s.store.TransferTx(ctx, db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	})
	if err != nil {
		return domain.TransferResponse{}, err
	}

	return domain.TransferResponse{
		Transfer:    toDomainTransfer(result.Transfer),
		FromAccount: toDomainAccount(result.FromAccount),
		ToAccount:   toDomainAccount(result.ToAccount),
		FromEntry:   toDomainEntry(result.FromEntry),
		ToEntry:     toDomainEntry(result.ToEntry),
	}, nil
}

func (s *transferService) GetTransfer(ctx context.Context, id int64) (domain.Transfer, error) {
	if id <= 0 {
		return domain.Transfer{ID: 0}, domain.ErrTransferNotFound
	}

	transfer, err := s.store.GetTransfer(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Transfer{}, domain.ErrTransferNotFound
		}
		return domain.Transfer{}, err
	}

	return toDomainTransfer(transfer), nil
}

func (s *transferService) ListTransfers(ctx context.Context, req domain.ListTransfersRequest) ([]domain.Transfer, error) {
	pageID := req.PageID
	if pageID <= 0 {
		pageID = 1
	}

	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	} else if pageSize > 1000 {
		pageSize = 1000
	}

	limit := pageSize
	offset := (pageID - 1) * pageSize

	transfers, err := s.store.ListTransfers(ctx, db.ListTransfersParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	return toDomainTransfers(transfers), nil
}
