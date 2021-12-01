package service

import (
	"time"

	"github.com/mstreet3/banking/domain"
	"github.com/mstreet3/banking/dto"
	"github.com/mstreet3/banking/errs"
)

type TransactionService interface {
	NewTransaction(dto.NewTransactionRequest) (*dto.NewTransactionResponse, *errs.AppError)
}

type DefaultTransactionService struct {
	repo domain.TransactionRepository
}

func NewTransactionService(repo domain.TransactionRepository) DefaultTransactionService {
	return DefaultTransactionService{repo}
}

func (s DefaultTransactionService) NewTransaction(req dto.NewTransactionRequest) (*dto.NewTransactionResponse, *errs.AppError) {
	invalid := req.Validate()
	if invalid != nil {
		return nil, invalid
	}
	t := domain.Transaction{
		AccountId:       req.AccountId,
		TransactionType: req.TransactionType,
		Amount:          req.Amount,
		TransactionDate: time.Now().Format("2006-01-02 15:04:05"),
	}
	newTxn, updAcct, err := s.repo.Save(t)
	if err != nil {
		return nil, err
	}
	response := newTxn.ToNewTransactionResponseDto()
	response.AccountBalance = updAcct.Amount
	return &response, nil
}
