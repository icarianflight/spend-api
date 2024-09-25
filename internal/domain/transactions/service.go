package transactions

import "time"

// TransactionService provides the core logic for managing transactions.
type TransactionService struct {
	transactionPersistence ForSavingTransaction
}

// NewTransactionService creates a new TransactionService.
func NewTransactionService(persistence ForSavingTransaction) *TransactionService {
	return &TransactionService{
		transactionPersistence: persistence,
	}
}

// CreateTransaction creates a new transaction and saves it using persistence.
func (s *TransactionService) CreateTransaction(accountID string, amount float64, txnType, description string) (*Transaction, error) {
	transaction := NewTransaction("", accountID, amount, txnType, time.Now(), description)

	// Save the transaction using the persistence port
	err := s.transactionPersistence.SaveTransaction(transaction)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}
