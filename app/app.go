package app

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/mstreet3/banking/domain"
	"github.com/mstreet3/banking/service"
)

func getDbClient() *sqlx.DB {
	client, err := sqlx.Open("mysql", "root:codecamp@tcp(localhost:3306)/banking")
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	client.SetConnMaxLifetime(time.Minute * 3)
	client.SetMaxOpenConns(10)
	client.SetMaxIdleConns(10)
	return client
}

func Start() {
	router := mux.NewRouter()

	/* define data source */
	dbClient := getDbClient()
	customerRepository := domain.NewCustomerRepositoryDb(dbClient)
	accountRepository := domain.NewAccountRepositoryDb(dbClient)

	/* create a customer service from the data source */
	customerService := service.NewCustomerService(customerRepository)
	accountService := service.NewAccountService(accountRepository)

	/* define the service used by the customer handlers */
	ch := CustomerHandlers{
		service: customerService,
	}

	ah := AccountHandlers{
		service: accountService,
	}

	/* implement a get all customers route */
	router.Path("/customers").Queries("status", "{status}").HandlerFunc(ch.getAllCustomersByStatus).Methods(http.MethodGet)
	router.Path("/customers").HandlerFunc(ch.getAllCustomers).Methods(http.MethodGet)
	router.HandleFunc("/customers/{customer_id:[0-9]+}", ch.getCustomer).Methods(http.MethodGet)
	router.HandleFunc("/customers/{customer_id:[0-9]+}/account", ah.newAccount).Methods(http.MethodPost)

	err := http.ListenAndServe("localhost:8080", router)

	if err != nil {
		log.Fatal(err)
	}
}
