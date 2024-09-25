package transactions

import "time"

// Transaction represents a financial transaction associated with an account.
type Transaction struct {
	ID          string
	AccountID   string
	Amount      float64
	Type        string
	Timestamp   time.Time
	Description string
}

// NewTransaction creates a new transaction.
func NewTransaction(id, accountID string, amount float64, txnType string, timestamp time.Time, description string) *Transaction {
	return &Transaction{
		ID:          id,
		AccountID:   accountID,
		Amount:      amount,
		Type:        txnType,
		Timestamp:   timestamp,
		Description: description,
	}
}
