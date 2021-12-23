package app

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mstreet3/banking/dto"
	"github.com/mstreet3/banking/service"
	"github.com/mstreet3/banking/utils"
)

type CustomerHandlers struct {
	service service.CustomerService
}

func (ch *CustomerHandlers) getAllCustomers(w http.ResponseWriter, r *http.Request) {
	/* access the customer repository via the service */
	customers, err := ch.service.GetAllCustomers()

	/* handle the response from the service */
	if err != nil {
		utils.WriteResponse(w, err.Code, err.AsMessage())
	} else {
		utils.WriteResponse(w, http.StatusOK, customers)
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
		utils.WriteResponse(w, err.Code, err.AsMessage())
	} else {
		utils.WriteResponse(w, http.StatusOK, customers)
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
		utils.WriteResponse(w, err.Code, err.AsMessage())
	} else {
		utils.WriteResponse(w, http.StatusOK, customer)
	}

}
