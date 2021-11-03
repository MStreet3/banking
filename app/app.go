package app

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mstreet3/banking/domain"
	"github.com/mstreet3/banking/service"
)

func Start() {
	router := mux.NewRouter()

	/* define data source */
	customerRepository := domain.NewCustomerRepositoryStub()

	/* create a customer service from the data source */
	customerService := service.NewCustomerService(customerRepository)

	/* define the service used by the customer handlers */
	ch := CustomerHandlers{
		service: customerService,
	}

	/* implement a get all customers route */
	router.HandleFunc("/customers", ch.getAllCustomers)

	err := http.ListenAndServe("localhost:8080", router)

	if err != nil {
		log.Fatal(err)
	}
}
