package accounts

// Account represents a bank account with an ID and a Name.
type Account struct {
	ID   string
	Name string
}

// NewAccount creates a new account with the given ID and Name.
func NewAccount(id, name string) *Account {
	return &Account{
		ID:   id,
		Name: name,
	}
}
