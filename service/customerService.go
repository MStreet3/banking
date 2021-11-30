package service

import (
	"github.com/mstreet3/banking/domain"
	"github.com/mstreet3/banking/errs"
)

type CustomerService interface {
	GetAllCustomers() ([]domain.Customer, *errs.AppError)
	GetAllCustomersByStatus(status domain.CustomerStatus) ([]domain.Customer, *errs.AppError)
	GetCustomer(string) (*domain.Customer, *errs.AppError)
}

type DefaultCustomerService struct {
	repo domain.CustomerRepository
}

func (s DefaultCustomerService) GetAllCustomers() ([]domain.Customer, *errs.AppError) {
	return s.repo.FindAll()
}

func (s DefaultCustomerService) GetAllCustomersByStatus(status domain.CustomerStatus) ([]domain.Customer, *errs.AppError) {
	return s.repo.FindAllByStatus(status)
}

func (s DefaultCustomerService) GetCustomer(id string) (*domain.Customer, *errs.AppError) {
	return s.repo.ById(id)
}

func NewCustomerService(repo domain.CustomerRepository) DefaultCustomerService {
	return DefaultCustomerService{
		repo: repo,
	}
}
