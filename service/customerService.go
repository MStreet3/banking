package service

import (
	"github.com/mstreet3/banking/domain"
	"github.com/mstreet3/banking/dto"
	"github.com/mstreet3/banking/errs"
)

type CustomerService interface {
	GetAllCustomers() ([]dto.CustomerResponse, *errs.AppError)
	GetAllCustomersByStatus(status dto.CustomerStatus) ([]dto.CustomerResponse, *errs.AppError)
	GetCustomer(string) (*dto.CustomerResponse, *errs.AppError)
}

type DefaultCustomerService struct {
	repo domain.CustomerRepository
}

func (s DefaultCustomerService) GetAllCustomers() ([]dto.CustomerResponse, *errs.AppError) {
	c, err := s.repo.FindAll()
	return handleCustomerFetch(c, err)
}

func (s DefaultCustomerService) GetAllCustomersByStatus(status dto.CustomerStatus) ([]dto.CustomerResponse, *errs.AppError) {
	c, err := s.repo.FindAllByStatus(status)
	return handleCustomerFetch(c, err)

}

func (s DefaultCustomerService) GetCustomer(id string) (*dto.CustomerResponse, *errs.AppError) {
	c, err := s.repo.ById(id)
	if err != nil {
		return nil, err
	}
	response := c.ToDto()
	return &response, nil
}

func NewCustomerService(repo domain.CustomerRepository) DefaultCustomerService {
	return DefaultCustomerService{
		repo: repo,
	}
}

func handleCustomerFetch(c []domain.Customer, err *errs.AppError) ([]dto.CustomerResponse, *errs.AppError) {
	if err != nil {
		return nil, err
	}
	response := make([]dto.CustomerResponse, 0)
	for _, cus := range c {
		response = append(response, cus.ToDto())
	}
	return response, nil
}
