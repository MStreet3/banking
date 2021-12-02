package domain

import (
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/mstreet3/banking/errs"
	"github.com/mstreet3/banking/logger"
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
		logger.Error("error while creating new account " + err.Error())
		return nil, errs.UnexpectedDatabaseError()
	}
	id, err := result.LastInsertId()
	if err != nil {
		logger.Error("error while getting id of new account " + err.Error())
		return nil, errs.UnexpectedDatabaseError()
	}
	a.AccountId = strconv.FormatInt(id, 10)
	return &a, nil
}

func (db AccountRepositoryDb) BalanceCheck(t Transaction) *errs.AppError {
	// Getting the latest account information from the accounts table
	acct, err := db.FindBy(t.AccountId)
	if err != nil {
		return err
	}

	if t.IsWithdrawal() && !acct.CanWithdraw(t.Amount) {
		return errs.InsufficientFundsError()
	}
	return nil
}

func (db AccountRepositoryDb) AddTransaction(t Transaction) (*Transaction, *Account, *errs.AppError) {

	// validate the transaction to add
	invalid := db.BalanceCheck(t)
	if invalid != nil {
		return nil, nil, invalid
	}

	// start the db transaction
	tx, err := db.client.Begin()
	if err != nil {
		logger.Error("Error while starting a new transaction for bank account transaction: " + err.Error())
		return nil, nil, errs.UnexpectedDatabaseError()
	}

	result, err := tx.Exec(`INSERT INTO transactions (account_id, amount, transaction_type, transaction_date) 
											values (?, ?, ?, ?)`, t.AccountId, t.Amount, t.TransactionType, t.TransactionDate)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while saving transaction: " + err.Error())
		return nil, nil, errs.UnexpectedDatabaseError()
	}

	// updating account balance
	if t.IsWithdrawal() {
		_, err = tx.Exec(`UPDATE accounts SET amount = amount - ? where account_id = ?`, t.Amount, t.AccountId)
	} else {
		_, err = tx.Exec(`UPDATE accounts SET amount = amount + ? where account_id = ?`, t.Amount, t.AccountId)
	}

	// in case of error Rollback, and changes from both the tables will be reverted
	if err != nil {
		tx.Rollback()
		logger.Error("Error while saving transaction: " + err.Error())
		return nil, nil, errs.UnexpectedDatabaseError()
	}
	// commit the transaction when all is good
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		logger.Error("Error while commiting transaction for bank account: " + err.Error())
		return nil, nil, errs.UnexpectedDatabaseError()
	}

	// getting the last transaction ID from the transaction table
	transactionId, err := result.LastInsertId()
	if err != nil {
		logger.Error("Error while getting the last transaction id: " + err.Error())
		return nil, nil, errs.UnexpectedDatabaseError()
	}

	// Getting the latest account information from the accounts table
	acct, appErr := db.FindBy(t.AccountId)
	if appErr != nil {
		return nil, nil, appErr
	}

	t.TransactionId = strconv.FormatInt(transactionId, 10)
	return &t, acct, nil
}

func (d AccountRepositoryDb) FindBy(accountId string) (*Account, *errs.AppError) {
	sqlGetAccount := "SELECT account_id, customer_id, opening_date, account_type, amount from accounts where account_id = ?"
	var account Account
	err := d.client.Get(&account, sqlGetAccount, accountId)
	if err != nil {
		logger.Error("Error while fetching account information: " + err.Error())
		return nil, errs.UnexpectedDatabaseError()
	}
	return &account, nil
}
