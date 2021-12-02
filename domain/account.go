package domain

import (
	"github.com/mstreet3/banking/dto"
	"github.com/mstreet3/banking/errs"
)

type AccountType string

const (
	SAVING   AccountType = "saving"
	CHECKING AccountType = "checking"
)

type Account struct {
	AccountId   string      `db:"account_id"`
	CustomerId  string      `db:"customer_id"`
	OpeningDate string      `db:"opening_date"`
	AccountType AccountType `db:"account_type"`
	Amount      float64
	Status      string
}

type AccountRepository interface {
	Save(Account) (*Account, *errs.AppError)
	AddTransaction(Transaction) (*Transaction, *Account, *errs.AppError)
}

func (a Account) ToNewAccountResponseDto() dto.NewAccountResponse {
	return dto.NewAccountResponse{AccountId: a.AccountId}
}

func (a Account) CanWithdraw(amount float64) bool {
	return a.Amount > amount
}
