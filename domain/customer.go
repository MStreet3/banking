package domain

import "github.com/mstreet3/banking/errs"

type CustomerStatus string

const (
	ACTIVE   CustomerStatus = "active"
	INACTIVE CustomerStatus = "inactive"
)

type Customer struct {
	Id          string `db:"customer_id"`
	Name        string
	City        string
	Zipcode     string
	DateOfBirth string `db:"date_of_birth"`
	Status      string
}

type CustomerRepository interface {
	FindAll() ([]Customer, *errs.AppError)
	FindAllByStatus(status CustomerStatus) ([]Customer, *errs.AppError)
	ById(string) (*Customer, *errs.AppError)
}
