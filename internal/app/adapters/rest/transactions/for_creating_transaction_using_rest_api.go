package transactions

import (
	"encoding/json"
	"net/http"
	"spend-api/internal/domain/transactions"
)

// ForCreatingTransactionUsingRestAPI is the REST API adapter for creating transactions.
type ForCreatingTransactionUsingRestAPI struct {
	transactionService transactions.ForCreatingTransaction
}

// NewForCreatingTransactionUsingRestAPI creates a new REST handler for creating transactions.
func NewForCreatingTransactionUsingRestAPI(service transactions.ForCreatingTransaction) *ForCreatingTransactionUsingRestAPI {
	return &ForCreatingTransactionUsingRestAPI{
		transactionService: service,
	}
}

// ServeHTTP handles HTTP requests for creating a transaction.
func (h *ForCreatingTransactionUsingRestAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	var requestBody struct {
		AccountID   string  `json:"accountID"`
		Amount      float64 `json:"amount"`
		Type        string  `json:"type"`
		Description string  `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	transaction, err := h.transactionService.CreateTransaction(requestBody.AccountID, requestBody.Amount, requestBody.Type, requestBody.Description)
	if err != nil {
		http.Error(w, "Failed to create transaction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(transaction)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
