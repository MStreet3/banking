package dto

import "github.com/mstreet3/banking/errs"

type NewTransactionRequest struct {
	AccountId       string `json:"account_id"`
	Amount          float64
	TransactionType string `json:"transaction_type"`
}

type NewTransactionResponse struct {
	TransactionId   string  `json:"transaction_id"`
	AccountId       string  `json:"account_id"`
	Amount          float64 `json:"new_balance"`
	TransactionType string  `json:"transaction_type"`
	TransactionDate string  `json:"transaction_date"`
}

func (req NewTransactionRequest) Validate() *errs.AppError {
	if req.TransactionType != "withdrawal" && req.TransactionType != "deposit" {
		return errs.NewValidationError("transaction type must either be withdrawal or deposit")
	}
	if req.Amount < 0 {
		return errs.NewValidationError("transaction value must be positive")
	}
	return nil
}
