package transactions

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"spend-api/internal/domain/transactions"
	"testing"
	"time"
)

// Test for successful transaction creation via the REST API
func TestForCreatingTransactionUsingRestAPI_Success(t *testing.T) {
	fakeTransactionService := &FakeForCreatingTransaction{}

	apiHandler := NewForCreatingTransactionUsingRestAPI(fakeTransactionService)

	requestBody := map[string]interface{}{
		"accountID":   "12345",
		"amount":      100.0,
		"type":        "credit",
		"description": "Payment",
	}
	jsonBody, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	respRecorder := httptest.NewRecorder()

	apiHandler.ServeHTTP(respRecorder, req)

	assert.Equal(t, http.StatusCreated, respRecorder.Code)
	assert.Contains(t, respRecorder.Body.String(), `"ID":"txn123"`)
	assert.Contains(t, respRecorder.Body.String(), `"AccountID":"12345"`)
}

// Test for invalid HTTP method
func TestForCreatingTransactionUsingRestAPI_InvalidMethod(t *testing.T) {
	fakeTransactionService := &FakeForCreatingTransaction{}
	apiHandler := NewForCreatingTransactionUsingRestAPI(fakeTransactionService)

	req := httptest.NewRequest(http.MethodGet, "/transactions", nil)
	respRecorder := httptest.NewRecorder()

	apiHandler.ServeHTTP(respRecorder, req)

	assert.Equal(t, http.StatusMethodNotAllowed, respRecorder.Code)
}

// Test for invalid JSON in the request body
func TestForCreatingTransactionUsingRestAPI_InvalidJSON(t *testing.T) {
	fakeTransactionService := &FakeForCreatingTransaction{}
	apiHandler := NewForCreatingTransactionUsingRestAPI(fakeTransactionService)

	req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader([]byte(`invalid json`)))
	req.Header.Set("Content-Type", "application/json")
	respRecorder := httptest.NewRecorder()

	apiHandler.ServeHTTP(respRecorder, req)

	assert.Equal(t, http.StatusBadRequest, respRecorder.Code)
}

// Test for internal server error from the transaction service
func TestForCreatingTransactionUsingRestAPI_ServiceError(t *testing.T) {
	fakeTransactionService := &FakeForCreatingTransaction{
		ReturnError: true,
	}
	apiHandler := NewForCreatingTransactionUsingRestAPI(fakeTransactionService)

	requestBody := map[string]interface{}{
		"accountID":   "12345",
		"amount":      100.0,
		"type":        "credit",
		"description": "Payment",
	}
	jsonBody, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	respRecorder := httptest.NewRecorder()

	apiHandler.ServeHTTP(respRecorder, req)

	assert.Equal(t, http.StatusInternalServerError, respRecorder.Code)
}

// Test for JSON encoding failure when creating a transaction
func TestForCreatingTransactionUsingRestAPI_EncodingError(t *testing.T) {
	fakeTransactionService := &FakeForCreatingTransaction{}
	apiHandler := NewForCreatingTransactionUsingRestAPI(fakeTransactionService)

	requestBody := map[string]interface{}{
		"accountID":   "12345",
		"amount":      100.0,
		"type":        "credit",
		"description": "Payment",
	}
	jsonBody, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	respRecorder := &errorResponseWriter{
		statusCode: http.StatusOK,
	}

	apiHandler.ServeHTTP(respRecorder, req)

	assert.Equal(t, http.StatusInternalServerError, respRecorder.statusCode)
}

// Fake ResponseWriter that simulates an encoding failure
type errorResponseWriter struct {
	statusCode int
}

func (e *errorResponseWriter) Header() http.Header {
	return http.Header{}
}

func (e *errorResponseWriter) WriteHeader(statusCode int) {
	e.statusCode = statusCode
}

func (e *errorResponseWriter) Write(b []byte) (int, error) {
	return 0, io.ErrClosedPipe
}

// FakeForCreatingTransaction simulates the transaction service for testing.
type FakeForCreatingTransaction struct {
	ReturnError bool
}

func (f *FakeForCreatingTransaction) CreateTransaction(accountID string, amount float64, txnType, description string) (*transactions.Transaction, error) {
	if f.ReturnError {
		return nil, errors.New("failed to create transaction")
	}
	return &transactions.Transaction{
		ID:          "txn123",
		AccountID:   accountID,
		Amount:      amount,
		Type:        txnType,
		Timestamp:   time.Now(),
		Description: description,
	}, nil
}
