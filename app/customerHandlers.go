package app

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mstreet3/banking/dto"
	"github.com/mstreet3/banking/service"
)

type CustomerHandlers struct {
	service service.CustomerService
}

func (ch *CustomerHandlers) getAllCustomers(w http.ResponseWriter, r *http.Request) {
	/* access the customer repository via the service */
	customers, err := ch.service.GetAllCustomers()

	/* handle the response from the service */
	if err != nil {
		writeResponse(w, err.Code, err.AsMessage())
	} else {
		writeResponse(w, http.StatusOK, customers)
	}
}

func (ch *CustomerHandlers) getAllCustomersByStatus(w http.ResponseWriter, r *http.Request) {
	/* fetch URL params */
	vars := mux.Vars(r)
	status := vars["status"]

	/* access the customer repository via the service */
	customers, err := ch.service.GetAllCustomersByStatus(dto.CustomerStatus(status))

	/* handle the response from the service */
	if err != nil {
		writeResponse(w, err.Code, err.AsMessage())
	} else {
		writeResponse(w, http.StatusOK, customers)
	}
}

func (ch *CustomerHandlers) getCustomer(w http.ResponseWriter, r *http.Request) {
	/* fetch URL params */
	vars := mux.Vars(r)
	id := vars["customer_id"]

	/* access the customer repository via the service */
	customer, err := ch.service.GetCustomer(id)

	/* handle the response from the service */
	if err != nil {
		writeResponse(w, err.Code, err.AsMessage())
	} else {
		writeResponse(w, http.StatusOK, customer)
	}

}

func writeResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}

}
