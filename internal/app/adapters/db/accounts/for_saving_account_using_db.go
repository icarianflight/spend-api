package accounts

import (
	"fmt"
	"spend-api/internal/domain/accounts"
	"spend-api/internal/infra/db"
)

// ForSavingAccountUsingDB is the adapter for saving accounts using DB
type ForSavingAccountUsingDB struct {
	db db.Executor
}

// NewForSavingAccountUsingDB creates a new DB adapter for saving accounts
func NewForSavingAccountUsingDB(db db.Executor) *ForSavingAccountUsingDB {
	return &ForSavingAccountUsingDB{db: db}
}

// SaveAccount saves the given account to DB
func (a *ForSavingAccountUsingDB) SaveAccount(account *accounts.Account) error {
	query := "INSERT INTO accounts (name) VALUES (?)"
	result, err := a.db.Exec(query, account.Name)
	if err != nil {
		return fmt.Errorf("failed to save account: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to retrieve last insert ID: %w", err)
	}
	account.ID = fmt.Sprintf("%d", id)
	return nil
}
