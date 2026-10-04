package handler

import (
	"net/http"
	"strconv"

	"github.com/AdhamHussam/wallet-ledge/internal/domain"
	"github.com/AdhamHussam/wallet-ledge/internal/service"
)

// AccountHandler exposes account operations over HTTP.
type AccountHandler struct {
	svc service.AccountService
}

// NewAccountHandler creates a new AccountHandler.
func NewAccountHandler(svc service.AccountService) *AccountHandler {
	return &AccountHandler{svc: svc}
}

// Create handles POST /accounts.
func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateAccountRequest
	if err := readJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	account, err := h.svc.CreateAccount(r.Context(), req)
	if err != nil {
		writeDomainError(w, r, err)
		return
	}

	writeJSON(w, http.StatusCreated, account)
}

// Get handles GET /accounts/{id}.
func (h *AccountHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := parseIDParam(w, r)
	if !ok {
		return
	}

	account, err := h.svc.GetAccount(r.Context(), id)
	if err != nil {
		writeDomainError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, account)
}

// List handles GET /accounts?page_id=1&page_size=10.
func (h *AccountHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	pageID, err := parseOptionalInt32(q.Get("page_id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "page_id must be a positive integer")
		return
	}

	pageSize, err := parseOptionalInt32(q.Get("page_size"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "page_size must be a positive integer")
		return
	}

	accounts, err := h.svc.ListAccounts(r.Context(), domain.ListAccountsRequest{
		PageID:   pageID,
		PageSize: pageSize,
	})
	if err != nil {
		writeDomainError(w, r, err)
		return
	}

	// Always return a JSON array, never null.
	if accounts == nil {
		accounts = []domain.Account{}
	}

	writeJSON(w, http.StatusOK, accounts)
}

// parseIDParam extracts and validates the {id} path parameter.
func parseIDParam(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "id must be a positive integer")
		return 0, false
	}
	return id, true
}

// parseOptionalInt32 parses a positive int32 query value; empty means 0 (use service default).
func parseOptionalInt32(raw string) (int32, error) {
	if raw == "" {
		return 0, nil
	}
	v, err := strconv.ParseInt(raw, 10, 32)
	if err != nil || v <= 0 {
		return 0, strconv.ErrSyntax
	}
	return int32(v), nil
}
