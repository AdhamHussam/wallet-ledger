package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AdhamHussam/wallet-ledge/internal/domain"
)

// ---- mocks ----

type mockAccountService struct {
	createFn func(ctx context.Context, req domain.CreateAccountRequest) (domain.Account, error)
	getFn    func(ctx context.Context, id int64) (domain.Account, error)
	listFn   func(ctx context.Context, req domain.ListAccountsRequest) ([]domain.Account, error)
}

func (m *mockAccountService) CreateAccount(ctx context.Context, req domain.CreateAccountRequest) (domain.Account, error) {
	return m.createFn(ctx, req)
}
func (m *mockAccountService) GetAccount(ctx context.Context, id int64) (domain.Account, error) {
	return m.getFn(ctx, id)
}
func (m *mockAccountService) ListAccounts(ctx context.Context, req domain.ListAccountsRequest) ([]domain.Account, error) {
	return m.listFn(ctx, req)
}

type mockTransferService struct {
	createFn func(ctx context.Context, req domain.TransferRequest) (domain.TransferResponse, error)
	getFn    func(ctx context.Context, id int64) (domain.Transfer, error)
}

func (m *mockTransferService) CreateTransfer(ctx context.Context, req domain.TransferRequest) (domain.TransferResponse, error) {
	return m.createFn(ctx, req)
}
func (m *mockTransferService) GetTransfer(ctx context.Context, id int64) (domain.Transfer, error) {
	return m.getFn(ctx, id)
}

// ---- helpers ----

func doRequest(t *testing.T, h http.Handler, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	var reader *bytes.Reader
	if body != "" {
		reader = bytes.NewReader([]byte(body))
	} else {
		reader = bytes.NewReader(nil)
	}
	req := httptest.NewRequest(method, path, reader)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func decodeError(t *testing.T, rec *httptest.ResponseRecorder) string {
	t.Helper()
	var body errorBody
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode error body %q: %v", rec.Body.String(), err)
	}
	return body.Error
}

var sampleAccount = domain.Account{
	ID: 1, Owner: "alice", Balance: 100, Currency: domain.USD, CreatedAt: time.Now(),
}

// ---- tests ----

func TestHealth(t *testing.T) {
	router := NewRouter(&mockAccountService{}, &mockTransferService{})
	rec := doRequest(t, router, http.MethodGet, "/health", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestCreateAccountHandler(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		svcErr     error
		wantStatus int
	}{
		{"success", `{"owner":"alice","currency":"USD"}`, nil, http.StatusCreated},
		{"malformed json", `{"owner":`, nil, http.StatusBadRequest},
		{"empty body", ``, nil, http.StatusBadRequest},
		{"unknown field", `{"owner":"a","currency":"USD","admin":true}`, nil, http.StatusBadRequest},
		{"wrong type", `{"owner":123,"currency":"USD"}`, nil, http.StatusBadRequest},
		{"trailing data", `{"owner":"a","currency":"USD"}{}`, nil, http.StatusBadRequest},
		{"invalid owner", `{"owner":"","currency":"USD"}`, domain.ErrInvalidOwner, http.StatusBadRequest},
		{"invalid currency", `{"owner":"a","currency":"XXX"}`, domain.ErrInvalidCurrency, http.StatusBadRequest},
		{"internal error", `{"owner":"a","currency":"USD"}`, errors.New("db down"), http.StatusInternalServerError},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := &mockAccountService{
				createFn: func(_ context.Context, req domain.CreateAccountRequest) (domain.Account, error) {
					if tc.svcErr != nil {
						return domain.Account{}, tc.svcErr
					}
					return sampleAccount, nil
				},
			}
			router := NewRouter(svc, &mockTransferService{})
			rec := doRequest(t, router, http.MethodPost, "/accounts", tc.body)

			if rec.Code != tc.wantStatus {
				t.Fatalf("expected %d, got %d (body: %s)", tc.wantStatus, rec.Code, rec.Body.String())
			}
			if tc.wantStatus == http.StatusInternalServerError {
				if msg := decodeError(t, rec); msg != "internal server error" {
					t.Errorf("internal error leaked details: %q", msg)
				}
			}
		})
	}
}

func TestGetAccountHandler(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		svcErr     error
		wantStatus int
	}{
		{"success", "/accounts/1", nil, http.StatusOK},
		{"not found", "/accounts/42", domain.ErrAccountNotFound, http.StatusNotFound},
		{"non-numeric id", "/accounts/abc", nil, http.StatusBadRequest},
		{"zero id", "/accounts/0", nil, http.StatusBadRequest},
		{"negative id", "/accounts/-5", nil, http.StatusBadRequest},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := &mockAccountService{
				getFn: func(_ context.Context, id int64) (domain.Account, error) {
					if tc.svcErr != nil {
						return domain.Account{}, tc.svcErr
					}
					return sampleAccount, nil
				},
			}
			router := NewRouter(svc, &mockTransferService{})
			rec := doRequest(t, router, http.MethodGet, tc.path, "")
			if rec.Code != tc.wantStatus {
				t.Fatalf("expected %d, got %d (body: %s)", tc.wantStatus, rec.Code, rec.Body.String())
			}
		})
	}
}

func TestListAccountsHandler(t *testing.T) {
	t.Run("passes pagination params", func(t *testing.T) {
		var got domain.ListAccountsRequest
		svc := &mockAccountService{
			listFn: func(_ context.Context, req domain.ListAccountsRequest) ([]domain.Account, error) {
				got = req
				return []domain.Account{sampleAccount}, nil
			},
		}
		router := NewRouter(svc, &mockTransferService{})
		rec := doRequest(t, router, http.MethodGet, "/accounts?page_id=2&page_size=5", "")
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}
		if got.PageID != 2 || got.PageSize != 5 {
			t.Errorf("unexpected pagination: %+v", got)
		}
	})

	t.Run("empty list returns [] not null", func(t *testing.T) {
		svc := &mockAccountService{
			listFn: func(_ context.Context, _ domain.ListAccountsRequest) ([]domain.Account, error) {
				return nil, nil
			},
		}
		router := NewRouter(svc, &mockTransferService{})
		rec := doRequest(t, router, http.MethodGet, "/accounts", "")
		if body := bytes.TrimSpace(rec.Body.Bytes()); string(body) != "[]" {
			t.Errorf("expected [], got %s", body)
		}
	})

	t.Run("invalid page_size", func(t *testing.T) {
		router := NewRouter(&mockAccountService{}, &mockTransferService{})
		rec := doRequest(t, router, http.MethodGet, "/accounts?page_size=-1", "")
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rec.Code)
		}
	})
}

func TestCreateTransferHandler(t *testing.T) {
	validBody := `{"from_account_id":1,"to_account_id":2,"amount":50,"currency":"USD"}`

	tests := []struct {
		name       string
		body       string
		svcErr     error
		wantStatus int
	}{
		{"success", validBody, nil, http.StatusCreated},
		{"malformed json", `{"amount":`, nil, http.StatusBadRequest},
		{"string amount", `{"from_account_id":1,"to_account_id":2,"amount":"50","currency":"USD"}`, nil, http.StatusBadRequest},
		{"invalid amount", validBody, domain.ErrInvalidAmount, http.StatusBadRequest},
		{"same account", validBody, domain.ErrSameAccountTransfer, http.StatusBadRequest},
		{"account not found", validBody, domain.ErrAccountNotFound, http.StatusNotFound},
		{"insufficient funds", validBody, domain.ErrInsufficientFunds, http.StatusUnprocessableEntity},
		{"currency mismatch", validBody, domain.ErrCurrencyMismatch, http.StatusUnprocessableEntity},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := &mockTransferService{
				createFn: func(_ context.Context, req domain.TransferRequest) (domain.TransferResponse, error) {
					if tc.svcErr != nil {
						return domain.TransferResponse{}, tc.svcErr
					}
					return domain.TransferResponse{
						Transfer: domain.Transfer{ID: 7, FromAccountID: req.FromAccountID, ToAccountID: req.ToAccountID, Amount: req.Amount},
					}, nil
				},
			}
			router := NewRouter(&mockAccountService{}, svc)
			rec := doRequest(t, router, http.MethodPost, "/transfers", tc.body)
			if rec.Code != tc.wantStatus {
				t.Fatalf("expected %d, got %d (body: %s)", tc.wantStatus, rec.Code, rec.Body.String())
			}
		})
	}
}

func TestGetTransferHandler(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		svcErr     error
		wantStatus int
	}{
		{"success", "/transfers/7", nil, http.StatusOK},
		{"not found", "/transfers/99", domain.ErrTransferNotFound, http.StatusNotFound},
		{"bad id", "/transfers/xyz", nil, http.StatusBadRequest},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := &mockTransferService{
				getFn: func(_ context.Context, id int64) (domain.Transfer, error) {
					if tc.svcErr != nil {
						return domain.Transfer{}, tc.svcErr
					}
					return domain.Transfer{ID: id}, nil
				},
			}
			router := NewRouter(&mockAccountService{}, svc)
			rec := doRequest(t, router, http.MethodGet, tc.path, "")
			if rec.Code != tc.wantStatus {
				t.Fatalf("expected %d, got %d (body: %s)", tc.wantStatus, rec.Code, rec.Body.String())
			}
		})
	}
}

func TestMethodNotAllowed(t *testing.T) {
	router := NewRouter(&mockAccountService{}, &mockTransferService{})
	rec := doRequest(t, router, http.MethodDelete, "/accounts/1", "")
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}

func TestRecovererCatchesPanic(t *testing.T) {
	svc := &mockAccountService{
		getFn: func(_ context.Context, _ int64) (domain.Account, error) {
			panic("boom")
		},
	}
	router := NewRouter(svc, &mockTransferService{})
	rec := doRequest(t, router, http.MethodGet, "/accounts/1", "")
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
}
