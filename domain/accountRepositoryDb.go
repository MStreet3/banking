package domain

import (
	"log"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/mstreet3/banking/errs"
)

type AccountRepositoryDb struct {
	client *sqlx.DB
}

func NewAccountRepositoryDb(dbClient *sqlx.DB) AccountRepositoryDb {
	return AccountRepositoryDb{
		client: dbClient,
	}
}

func (db AccountRepositoryDb) Save(a Account) (*Account, *errs.AppError) {
	insertAcctSql := "INSERT INTO accounts (customer_id, opening_date, account_type, amount, status) VALUES (?,?,?,?,?)"
	result, err := db.client.Exec(insertAcctSql, a.CustomerId, a.OpeningDate, a.AccountType, a.Amount, a.Status)
	if err != nil {
		log.Println("error while creating new account " + err.Error())
		return nil, errs.NewInternalServerError("Unexpected database error")
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Println("error while getting id of new account " + err.Error())
		return nil, errs.NewInternalServerError("Unexpected database error")
	}
	a.AccountId = strconv.FormatInt(id, 10)
	return &a, nil
}
