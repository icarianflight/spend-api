package accounts

import (
	"database/sql"
	"errors"
	"spend-api/internal/domain/accounts"
	"testing"

	"github.com/stretchr/testify/assert"
)

// FakeDB for simulating DB behavior
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

// Test successful account saving
func TestForSavingAccountUsingDB_Success(t *testing.T) {
	fakeDB := &FakeDB{}
	adapter := NewForSavingAccountUsingDB(fakeDB)

	account := &accounts.Account{
		ID:   "12345",
		Name: "John Doe",
	}

	err := adapter.SaveAccount(account)
	assert.Nil(t, err, "Expected no error when saving account")
}

// Test account saving failure
func TestForSavingAccountUsingDB_Failure(t *testing.T) {
	mockDB := &FakeDB{ReturnError: true, ReturnInsertError: false}
	adapter := NewForSavingAccountUsingDB(mockDB)

	account := &accounts.Account{
		ID:   "12345",
		Name: "John Doe",
	}

	err := adapter.SaveAccount(account)
	assert.NotNil(t, err, "Expected an error when saving account")
	assert.Equal(t, "failed to save account: failed to execute query", err.Error(), "Expected error message to match")
}

// Test last insert ID failure
func TestForSavingAccountUsingDB_InsertIdFailure(t *testing.T) {
	mockDB := &FakeDB{ReturnError: false, ReturnInsertError: true}
	adapter := NewForSavingAccountUsingDB(mockDB)

	account := &accounts.Account{
		ID:   "12345",
		Name: "John Doe",
	}

	err := adapter.SaveAccount(account)
	assert.NotNil(t, err, "Expected an error when saving account")
	assert.Equal(t, "failed to retrieve last insert ID: failed to execute query", err.Error(), "Expected error message to match")
}
