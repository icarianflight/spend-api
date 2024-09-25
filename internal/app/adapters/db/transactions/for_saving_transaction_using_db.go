package transactions

import (
	"fmt"
	"spend-api/internal/domain/transactions"
	"spend-api/internal/infra/db"
)

// ForSavingTransactionUsingDB is the adapter for saving transactions using DB
type ForSavingTransactionUsingDB struct {
	db db.Executor
}

// NewForSavingTransactionUsingDB creates a new DB adapter for saving transactions
func NewForSavingTransactionUsingDB(executor db.Executor) *ForSavingTransactionUsingDB {
	return &ForSavingTransactionUsingDB{db: executor}
}

// SaveTransaction saves the given transaction to DB
func (a *ForSavingTransactionUsingDB) SaveTransaction(transaction *transactions.Transaction) error {
	query := "INSERT INTO transactions (account_id, amount, type, description) VALUES (?, ?, ?, ?, ?)"
	result, err := a.db.Exec(query, transaction.AccountID, transaction.Amount, transaction.Type, transaction.Description)
	if err != nil {
		return fmt.Errorf("failed to save transaction: %w", err)
	}
	// Retrieve the last insert ID and assign it to the account's ID field
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to retrieve last insert ID: %w", err)
	}
	transaction.ID = fmt.Sprintf("%d", id)
	return nil
}
