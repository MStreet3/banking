package domain

import (
	"github.com/mstreet3/banking/dto"
)

type Transaction struct {
	TransactionId   string `db:"transaction_id"`
	AccountId       string `db:"account_id"`
	Amount          float64
	TransactionType string `db:"transaction_type"`
	TransactionDate string `db:"transaction_date"`
}

func (t Transaction) ToTransactionResponseDto() dto.NewTransactionResponse {
	return dto.NewTransactionResponse{
		TransactionId:   t.TransactionId,
		TransactionType: t.TransactionType,
		TransactionDate: t.TransactionDate,
	}
}

func (t Transaction) IsWithdrawal() bool {
	return t.TransactionType == "withdrawal"
}
