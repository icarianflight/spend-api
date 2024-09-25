package transactions

// ForCreatingTransaction defines the port for creating a transaction.
type ForCreatingTransaction interface {
	CreateTransaction(accountID string, amount float64, txnType, description string) (*Transaction, error)
}

// ForSavingTransaction defines the port for saving a transaction in the persistence layer.
type ForSavingTransaction interface {
	SaveTransaction(transaction *Transaction) error
}
