package accounts

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"spend-api/internal/domain/accounts"
	"testing"

	"github.com/stretchr/testify/assert"
)

// FakeForCreatingAccount simulates the account service for testing.
type FakeForCreatingAccount struct {
	ReturnError bool
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

func (f *FakeForCreatingAccount) CreateAccount(name string) (*accounts.Account, error) {
	if f.ReturnError {
		return nil, errors.New("failed to create account")
	}
	return &accounts.Account{
		ID:   "12345",
		Name: name,
	}, nil
}

// Test for creating an account via the REST API
func TestForCreatingAccountUsingRestAPI(t *testing.T) {

	fakeAccountService := &FakeForCreatingAccount{}

	apiHandler := NewForCreatingAccountUsingRestAPI(fakeAccountService)

	requestBody := map[string]string{
		"name": "John Doe",
	}
	jsonBody, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	respRecorder := httptest.NewRecorder()

	apiHandler.ServeHTTP(respRecorder, req)

	assert.Equal(t, http.StatusCreated, respRecorder.Code, "Expected HTTP 201 Created")
	assert.Contains(t, respRecorder.Body.String(), `"id":"12345"`, "Response should contain account ID")
}

// Test for invalid HTTP method
func TestForCreatingAccountUsingRestAPI_InvalidMethod(t *testing.T) {
	fakeAccountService := &FakeForCreatingAccount{}
	apiHandler := NewForCreatingAccountUsingRestAPI(fakeAccountService)

	req := httptest.NewRequest(http.MethodGet, "/accounts", nil)
	respRecorder := httptest.NewRecorder()

	apiHandler.ServeHTTP(respRecorder, req)

	assert.Equal(t, http.StatusMethodNotAllowed, respRecorder.Code)
}

// Test for invalid JSON in the request body
func TestForCreatingAccountUsingRestAPI_InvalidJSON(t *testing.T) {
	fakeAccountService := &FakeForCreatingAccount{}
	apiHandler := NewForCreatingAccountUsingRestAPI(fakeAccountService)

	req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewReader([]byte(`invalid json`)))
	req.Header.Set("Content-Type", "application/json")
	respRecorder := httptest.NewRecorder()

	apiHandler.ServeHTTP(respRecorder, req)

	assert.Equal(t, http.StatusBadRequest, respRecorder.Code)
}

// Test for internal server error from the account service
func TestForCreatingAccountUsingRestAPI_ServiceError(t *testing.T) {
	fakeAccountService := &FakeForCreatingAccount{
		ReturnError: true,
	}
	apiHandler := NewForCreatingAccountUsingRestAPI(fakeAccountService)

	requestBody := map[string]string{
		"name": "John Doe",
	}
	jsonBody, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	respRecorder := httptest.NewRecorder()

	apiHandler.ServeHTTP(respRecorder, req)

	assert.Equal(t, http.StatusInternalServerError, respRecorder.Code)
}

// Test for JSON encoding failure when creating an account
func TestForCreatingAccountUsingRestAPI_EncodingError(t *testing.T) {
	fakeAccountService := &FakeForCreatingAccount{}
	apiHandler := NewForCreatingAccountUsingRestAPI(fakeAccountService)

	requestBody := map[string]string{
		"name": "John Doe",
	}
	jsonBody, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	respRecorder := &errorResponseWriter{}

	apiHandler.ServeHTTP(respRecorder, req)

	assert.Equal(t, http.StatusInternalServerError, respRecorder.statusCode,
		"Expected Internal Server Error if JSON encoding fails")
}
