package accounts

import (
	"encoding/json"
	"net/http"
	"spend-api/internal/domain/accounts"
)

type ForCreatingAccountUsingRestAPI struct {
	accountService accounts.ForCreatingAccount
}

// NewForCreatingAccountUsingRestAPI creates a new REST handler for creating accounts.
func NewForCreatingAccountUsingRestAPI(service accounts.ForCreatingAccount) *ForCreatingAccountUsingRestAPI {
	return &ForCreatingAccountUsingRestAPI{
		accountService: service,
	}
}

// ServeHTTP handles HTTP requests for creating an account.
func (h *ForCreatingAccountUsingRestAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	var requestBody struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	account, err := h.accountService.CreateAccount(requestBody.Name)
	if err != nil {
		http.Error(w, "Failed to create account", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(map[string]string{"id": account.ID})
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
