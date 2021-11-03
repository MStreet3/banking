package app

import (
	"encoding/json"
	"net/http"

	"github.com/mstreet3/banking/service"
)

type CustomerHandlers struct {
	service service.CustomerService
}

func (ch *CustomerHandlers) getAllCustomers(w http.ResponseWriter, r *http.Request) {

	customers, _ := ch.service.GetAllCustomers()

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customers)

}
