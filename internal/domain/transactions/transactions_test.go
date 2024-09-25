package transactions

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// FakeForSavingTransaction simulates the persistence layer for testing.
type FakeForSavingTransaction struct {
	ReturnError bool
}

func (f *FakeForSavingTransaction) SaveTransaction(transaction *Transaction) error {
	if f.ReturnError {
		return errors.New("failed to save transaction")
	}
	return nil
}

// Test for creating a new transaction with AccountID and Description
func TestCreateTransaction(t *testing.T) {
	transactionID := "txn123"
	accountID := "12345"
	amount := 100.0
	txnType := "credit"
	timestamp := time.Now()
	description := "Payment for groceries"

	newTransaction := NewTransaction(transactionID, accountID, amount, txnType, timestamp, description)

	assert.Equal(t, transactionID, newTransaction.ID, "Transaction ID should be correctly set")
	assert.Equal(t, accountID, newTransaction.AccountID, "Transaction should be tied to the correct Account ID")
	assert.Equal(t, amount, newTransaction.Amount, "Transaction amount should be correctly set")
	assert.Equal(t, txnType, newTransaction.Type, "Transaction type should be correctly set")
	assert.Equal(t, timestamp, newTransaction.Timestamp, "Transaction timestamp should be correctly set")
	assert.Equal(t, description, newTransaction.Description, "Transaction description should be correctly set")
}

// Test for creating and saving a transaction using FakeTransactionPersistence
func TestTransactionServiceCreateTransaction(t *testing.T) {
	fakePersistence := &FakeForSavingTransaction{}
	transactionService := NewTransactionService(fakePersistence)

	accountID := "12345"
	amount := 100.0
	txnType := "credit"
	description := "Payment for groceries"

	newTransaction, err := transactionService.CreateTransaction(accountID, amount, txnType, description)

	assert.Nil(t, err, "Error should be nil when creating a transaction")
	assert.Equal(t, "", newTransaction.ID, "Transaction ID should be blank")
	assert.Equal(t, accountID, newTransaction.AccountID, "Account ID should be correctly set")
	assert.Equal(t, amount, newTransaction.Amount, "Transaction amount should be correctly set")
	assert.Equal(t, txnType, newTransaction.Type, "Transaction type should be correctly set")
	assert.Equal(t, description, newTransaction.Description, "Transaction description should be correctly set")
	assert.WithinDuration(t, time.Now(), newTransaction.Timestamp, time.Second, "Transaction timestamp should be close to the current time")
}

// Test transaction creation failure due to SaveTransaction error
func TestTransactionServiceCreateTransaction_SaveError(t *testing.T) {
	fakePersistence := &FakeForSavingTransaction{
		ReturnError: true,
	}
	transactionService := NewTransactionService(fakePersistence)

	accountID := "12345"
	amount := 100.0
	txnType := "credit"
	description := "Payment"
	newTransaction, err := transactionService.CreateTransaction(accountID, amount, txnType, description)

	assert.NotNil(t, err, "Expected an error when saving transaction")
	assert.Nil(t, newTransaction, "No transaction should be returned when there's a saving error")
}
