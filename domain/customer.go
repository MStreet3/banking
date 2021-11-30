package domain

import (
	"github.com/mstreet3/banking/dto"
	"github.com/mstreet3/banking/errs"
)

type Customer struct {
	Id          string `db:"customer_id"`
	Name        string
	City        string
	Zipcode     string
	DateOfBirth string `db:"date_of_birth"`
	Status      string
}

func (c Customer) statusAsText() dto.CustomerStatus {
	statusAsText := dto.ACTIVE
	if c.Status == "0" {
		statusAsText = dto.INACTIVE
	}
	return statusAsText
}

func (c Customer) ToDto() dto.CustomerResponse {
	return dto.CustomerResponse{
		Id:          c.Id,
		Name:        c.Name,
		City:        c.City,
		Zipcode:     c.Zipcode,
		DateOfBirth: c.DateOfBirth,
		Status:      c.statusAsText(),
	}
}

type CustomerRepository interface {
	FindAll() ([]Customer, *errs.AppError)
	FindAllByStatus(status dto.CustomerStatus) ([]Customer, *errs.AppError)
	ById(string) (*Customer, *errs.AppError)
}
