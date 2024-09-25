package accounts

// AccountService provides the core logic for managing accounts.
type AccountService struct {
	accountPersistence ForSavingAccount
}

// NewAccountService creates a new AccountService.
func NewAccountService(persistence ForSavingAccount) *AccountService {
	return &AccountService{
		accountPersistence: persistence,
	}
}

// CreateAccount creates a new account and saves it using persistence.
func (s *AccountService) CreateAccount(name string) (*Account, error) {
	account := &Account{
		Name: name,
	}

	// Save the account using the persistence port
	err := s.accountPersistence.SaveAccount(account)
	if err != nil {
		return nil, err
	}

	return account, nil
}
