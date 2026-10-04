package handler

import (
	"net/http"

	"github.com/AdhamHussam/wallet-ledger/internal/domain"
	"github.com/AdhamHussam/wallet-ledger/internal/service"
)

// TransferHandler exposes transfer operations over HTTP.
type TransferHandler struct {
	svc service.TransferService
}

// NewTransferHandler creates a new TransferHandler.
func NewTransferHandler(svc service.TransferService) *TransferHandler {
	return &TransferHandler{svc: svc}
}

// Create handles POST /transfers.
func (h *TransferHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req domain.TransferRequest
	if err := readJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.svc.CreateTransfer(r.Context(), req)
	if err != nil {
		writeDomainError(w, r, err)
		return
	}

	writeJSON(w, http.StatusCreated, result)
}

// Get handles GET /transfers/{id}.
func (h *TransferHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := parseIDParam(w, r)
	if !ok {
		return
	}

	transfer, err := h.svc.GetTransfer(r.Context(), id)
	if err != nil {
		writeDomainError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, transfer)
}

// List handles GET /transfers?page_id=1&page_size=20.
func (h *TransferHandler) List(w http.ResponseWriter, r *http.Request) {
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

	transfers, err := h.svc.ListTransfers(r.Context(), domain.ListTransfersRequest{
		PageID:   pageID,
		PageSize: pageSize,
	})
	if err != nil {
		writeDomainError(w, r, err)
		return
	}

	if transfers == nil {
		transfers = []domain.Transfer{}
	}

	writeJSON(w, http.StatusOK, transfers)
}
