package transactions

import (
	"database/sql"
	"errors"
	"spend-api/internal/domain/transactions"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Fake DB for simulating DB behavior
type FakeDB struct {
	ReturnError       bool
	ReturnInsertError bool
}

func (f *FakeDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	if f.ReturnError {
		return nil, errors.New("failed to execute query")
	}
	if f.ReturnInsertError {
		return &MockFailedResult{}, nil
	} else {
		return &MockResult{}, nil
	}
}

func (f *FakeDB) Close() error {
	return nil
}

type MockResult struct{}
type MockFailedResult struct{}

func (r *MockResult) LastInsertId() (int64, error) { return 0, nil }
func (r *MockResult) RowsAffected() (int64, error) { return 1, nil }

func (r *MockFailedResult) LastInsertId() (int64, error) {
	return 1, errors.New("failed to execute query")
}
func (r *MockFailedResult) RowsAffected() (int64, error) { return 1, nil }

// Test successful transaction saving
func TestForSavingTransactionUsingDB_Success(t *testing.T) {
	fakeDB := &FakeDB{}
	adapter := NewForSavingTransactionUsingDB(fakeDB)

	transaction := &transactions.Transaction{
		ID:          "txn123",
		AccountID:   "12345",
		Amount:      100.0,
		Type:        "credit",
		Description: "Payment",
	}

	err := adapter.SaveTransaction(transaction)
	assert.Nil(t, err, "Expected no error when saving transaction")
}

// Test transaction saving failure
func TestForSavingTransactionUsingDB_Failure(t *testing.T) {
	mockDB := &FakeDB{ReturnError: true}
	adapter := NewForSavingTransactionUsingDB(mockDB)

	transaction := &transactions.Transaction{
		ID:          "txn123",
		AccountID:   "12345",
		Amount:      100.0,
		Type:        "credit",
		Description: "Payment",
	}

	err := adapter.SaveTransaction(transaction)
	assert.NotNil(t, err, "Expected an error when saving transaction")
	assert.Equal(t, "failed to save transaction: failed to execute query", err.Error(), "Expected error message to match")
}

func TestForSavingTransactionUsingDB_InsertIdFailure(t *testing.T) {
	mockDB := &FakeDB{ReturnError: false, ReturnInsertError: true}
	adapter := NewForSavingTransactionUsingDB(mockDB)

	transaction := &transactions.Transaction{
		ID:          "txn123",
		AccountID:   "12345",
		Amount:      100.0,
		Type:        "credit",
		Description: "Payment",
	}

	err := adapter.SaveTransaction(transaction)
	assert.NotNil(t, err, "Expected an error when saving account")
	assert.Equal(t, "failed to retrieve last insert ID: failed to execute query", err.Error(), "Expected error message to match")
}
