package accounts

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

// FakeForSavingAccount simulates the persistence layer for testing.
type FakeForSavingAccount struct {
	ReturnError bool
}

func (f *FakeForSavingAccount) SaveAccount(account *Account) error {
	if f.ReturnError {
		return errors.New("failed to save account")
	}
	return nil
}

// Test for creating a new account
func TestCreateAccount(t *testing.T) {
	accountID := "12345"
	accountName := "An Account"
	newAccount := NewAccount(accountID, accountName)

	assert.Equal(t, accountID, newAccount.ID, "Account ID should be correctly set")
	assert.Equal(t, accountName, newAccount.Name, "New account should have a name set")
}

// Test for saving an account using persistence
func TestAccountServiceCreateAccount(t *testing.T) {
	fakePersistence := &FakeForSavingAccount{}
	accountService := NewAccountService(fakePersistence)

	accountName := "John Doe"
	newAccount, err := accountService.CreateAccount(accountName)

	assert.Nil(t, err, "Error should be nil when creating an account")
	assert.Equal(t, "", newAccount.ID, "Created account ID should be blank")
	assert.Equal(t, accountName, newAccount.Name, "Created account name should match")
}

// Test account creation failure due to SaveAccount error
func TestAccountServiceCreateAccount_SaveError(t *testing.T) {
	fakePersistence := &FakeForSavingAccount{
		ReturnError: true,
	}
	accountService := NewAccountService(fakePersistence)

	accountName := "John Doe"
	newAccount, err := accountService.CreateAccount(accountName)

	assert.NotNil(t, err, "Expected an error when saving account")
	assert.Nil(t, newAccount, "No account should be returned when there's a saving error")
}
