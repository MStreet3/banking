package domain

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/mstreet3/banking/errs"
)

type CustomerRepositoryDb struct {
	client *sqlx.DB
}

func NewCustomerRepositoryDb() CustomerRepositoryDb {
	client, err := sqlx.Open("mysql", "root:codecamp@tcp(localhost:3306)/banking")
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	client.SetConnMaxLifetime(time.Minute * 3)
	client.SetMaxOpenConns(10)
	client.SetMaxIdleConns(10)
	return CustomerRepositoryDb{
		client: client,
	}
}

func (d CustomerRepositoryDb) FindAllByStatus(status CustomerStatus) ([]Customer, *errs.AppError) {
	findAllSql := "select customer_id, name, city, zipcode, date_of_birth, status from customers where status = ?"
	var s int
	switch status {
	case ACTIVE:
		s = 1
	case INACTIVE:
		s = 0
	}
	return handleCustomerSelectQuery(d.client, findAllSql, s)
}

func (d CustomerRepositoryDb) FindAll() ([]Customer, *errs.AppError) {
	findAllSql := "select customer_id, name, city, zipcode, date_of_birth, status from customers"
	return handleCustomerSelectQuery(d.client, findAllSql)
}

func (d CustomerRepositoryDb) ById(id string) (*Customer, *errs.AppError) {
	findByIdSql := "select customer_id, name, city, zipcode, date_of_birth, status from customers where customer_id=?"
	var c Customer
	err := d.client.Get(&c, findByIdSql, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("Customer not found")
		}
		return nil, errs.NewInternalServerError("Unexpected database error")
	}
	return &c, nil
}

func handleCustomerSelectQuery(db *sqlx.DB, q string, args ...interface{}) ([]Customer, *errs.AppError) {
	customers := make([]Customer, 0)
	err := db.Select(&customers, q, args...)
	if err != nil {
		log.Println("error while scanning rows " + err.Error())
		return nil, errs.NewInternalServerError("Unexpected database error")
	}
	return customers, nil

}
