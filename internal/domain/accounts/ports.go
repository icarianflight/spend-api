package accounts

// ForCreatingAccount defines the port for creating an account.
type ForCreatingAccount interface {
	CreateAccount(name string) (*Account, error)
}

// ForSavingAccount defines the port for saving an account to persistence
type ForSavingAccount interface {
	SaveAccount(account *Account) error
}
