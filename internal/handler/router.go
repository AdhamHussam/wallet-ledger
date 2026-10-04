package handler

import (
	"net/http"

	"github.com/AdhamHussam/wallet-ledger/internal/service"
)

// NewRouter wires all HTTP routes and middleware.
func NewRouter(accounts service.AccountService, transfers service.TransferService) http.Handler {
	ah := NewAccountHandler(accounts)
	th := NewTransferHandler(transfers)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{
			"status":  "up",
			"service": "wallet-ledger",
		})
	})

	mux.HandleFunc("POST /accounts", ah.Create)
	mux.HandleFunc("GET /accounts", ah.List)
	mux.HandleFunc("GET /accounts/{id}", ah.Get)

	mux.HandleFunc("POST /transfers", th.Create)
	mux.HandleFunc("GET /transfers/{id}", th.Get)

	// Middleware order: CORS -> Logging -> Recoverer -> Handlers
	return CORS(Logging(Recoverer(mux)))
}
