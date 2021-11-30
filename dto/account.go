package dto

import "github.com/mstreet3/banking/errs"

type NewAccountRequest struct {
	CustomerId  string `json:"customer_id"`
	Amount      float64
	AccountType string `json:"account_type"`
}

type NewAccountResponse struct {
	AccountId string `json:"account_id"`
}

func (req NewAccountRequest) Validate() *errs.AppError {
	if req.AccountType != "checking" && req.AccountType != "saving" {
		return errs.NewValidationError("account type must either be checking or saving")
	}
	if req.Amount < 5000.00 {
		return errs.NewValidationError("minimum opening account balance is 5000.00")
	}
	return nil
}
