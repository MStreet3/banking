package service

import (
	"time"

	"github.com/mstreet3/banking/domain"
	"github.com/mstreet3/banking/dto"
	"github.com/mstreet3/banking/errs"
)

type AccountService interface {
	NewAccount(dto.NewAccountRequest) (*dto.NewAccountResponse, *errs.AppError)
}

type DefaultAccountService struct {
	repo domain.AccountRepository
}

func (s DefaultAccountService) NewAccount(req dto.NewAccountRequest) (*dto.NewAccountResponse, *errs.AppError) {
	invalid := req.Validate()
	if invalid != nil {
		return nil, invalid
	}
	a := domain.Account{
		CustomerId:  req.CustomerId,
		AccountType: domain.AccountType(req.AccountType),
		Amount:      req.Amount,
		OpeningDate: time.Now().Format("2006-01-02 15:04:05"),
		Status:      "1",
	}
	newAcct, err := s.repo.Save(a)
	if err != nil {
		return nil, err
	}
	response := newAcct.ToNewAccountResponseDto()
	return &response, nil
}

func NewAccountService(repo domain.AccountRepository) DefaultAccountService {
	return DefaultAccountService{repo}
}
