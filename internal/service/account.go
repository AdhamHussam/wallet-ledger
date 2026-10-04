package service

import (
	"context"
	"errors"
	"strings"

	"github.com/AdhamHussam/wallet-ledge/internal/domain"
	"github.com/AdhamHussam/wallet-ledge/internal/repository/db"
	"github.com/jackc/pgx/v5"
)

// AccountService defines business operations on accounts.
type AccountService interface {
	CreateAccount(ctx context.Context, req domain.CreateAccountRequest) (domain.Account, error)
	GetAccount(ctx context.Context, id int64) (domain.Account, error)
	ListAccounts(ctx context.Context, req domain.ListAccountsRequest) ([]domain.Account, error)
}

type accountService struct {
	store db.Store
}

// NewAccountService creates a new AccountService.
func NewAccountService(store db.Store) AccountService {
	return &accountService{store: store}
}

func (s *accountService) CreateAccount(ctx context.Context, req domain.CreateAccountRequest) (domain.Account, error) {
	owner := strings.TrimSpace(req.Owner)
	if owner == "" {
		return domain.Account{}, domain.ErrInvalidOwner
	}

	currency := strings.ToUpper(strings.TrimSpace(req.Currency))
	if !domain.IsSupportedCurrency(currency) {
		return domain.Account{}, domain.ErrInvalidCurrency
	}

	if req.InitialBalance < 0 {
		return domain.Account{}, domain.ErrInvalidAmount
	}

	account, err := s.store.CreateAccount(ctx, db.CreateAccountParams{
		Owner:    owner,
		Balance:  req.InitialBalance,
		Currency: currency,
	})
	if err != nil {
		return domain.Account{}, err
	}

	return toDomainAccount(account), nil
}

func (s *accountService) GetAccount(ctx context.Context, id int64) (domain.Account, error) {
	if id <= 0 {
		return domain.Account{}, domain.ErrAccountNotFound
	}

	account, err := s.store.GetAccount(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Account{}, domain.ErrAccountNotFound
		}
		return domain.Account{}, err
	}

	return toDomainAccount(account), nil
}

func (s *accountService) ListAccounts(ctx context.Context, req domain.ListAccountsRequest) ([]domain.Account, error) {
	pageID := req.PageID
	if pageID <= 0 {
		pageID = 1
	}

	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	} else if pageSize > 100 {
		pageSize = 100
	}

	limit := pageSize
	offset := (pageID - 1) * pageSize

	accounts, err := s.store.ListAccounts(ctx, db.ListAccountsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	return toDomainAccounts(accounts), nil
}
