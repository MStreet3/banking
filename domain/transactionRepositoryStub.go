package domain

import "github.com/mstreet3/banking/errs"

type TransactionRepositoryStub struct{}

func (s TransactionRepositoryStub) Save(t Transaction) (*Transaction, *Account, *errs.AppError) {
	txn := Transaction{
		TransactionId:   "10001",
		AccountId:       "2000",
		Amount:          1000.00,
		TransactionType: "deposit",
		TransactionDate: "2021-12-01 08:34",
	}
	acct := Account{
		AccountId: "2000",
		Amount:    6000.00,
	}
	return &txn, &acct, nil
}
