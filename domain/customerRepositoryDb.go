package domain

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mstreet3/banking/errs"
)

type CustomerRepositoryDb struct {
	client *sql.DB
}

func NewCustomerRepositoryDb() CustomerRepositoryDb {
	client, err := sql.Open("mysql", "root:codecamp@tcp(localhost:3306)/banking")
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
	rows, err := d.client.Query(findAllSql, s)
	return handleCustomerRows(rows, err)
}

func (d CustomerRepositoryDb) FindAll() ([]Customer, *errs.AppError) {
	findAllSql := "select customer_id, name, city, zipcode, date_of_birth, status from customers"
	rows, err := d.client.Query(findAllSql)
	return handleCustomerRows(rows, err)
}

func (d CustomerRepositoryDb) ById(id string) (*Customer, *errs.AppError) {
	findByIdSql := "select customer_id, name, city, zipcode, date_of_birth, status from customers where customer_id=?"
	row := d.client.QueryRow(findByIdSql, id)
	var c Customer
	err := row.Scan(&c.Id, &c.Name, &c.City, &c.Zipcode, &c.DateOfBirth, &c.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("Customer not found")
		}
		return nil, errs.NewInternalServerError("Unexpected database error")
	}
	return &c, nil
}

func handleCustomerRows(rows *sql.Rows, err error) ([]Customer, *errs.AppError) {
	if err != nil {
		log.Println(("error while querying customer table " + err.Error()))
		return nil, errs.NewInternalServerError("Unexpected database error")
	}
	customers := make([]Customer, 0)
	for rows.Next() {
		var c Customer
		err := rows.Scan(&c.Id, &c.Name, &c.City, &c.Zipcode, &c.DateOfBirth, &c.Status)
		if err != nil {
			log.Println("error while scanning rows " + err.Error())
			return nil, errs.NewInternalServerError("Unexpected database error")
		}
		customers = append(customers, c)
	}
	return customers, nil
}
