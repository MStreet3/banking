package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mstreet3/banking/service"
)

type CustomerHandlers struct {
	service service.CustomerService
}

func (ch *CustomerHandlers) getAllCustomers(w http.ResponseWriter, r *http.Request) {
	/* access the customer repository via the service */
	customers, _ := ch.service.GetAllCustomers()

	/* handle the response from the service */
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customers)

}

func (ch *CustomerHandlers) getCustomer(w http.ResponseWriter, r *http.Request) {
	/* fetch URL params */
	vars := mux.Vars(r)
	id := vars["customer_id"]

	/* access the customer repository via the service */
	customer, err := ch.service.GetCustomer(id)

	/* handle the response from the service */
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, err.Error())
	} else {
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(customer)
	}

}
