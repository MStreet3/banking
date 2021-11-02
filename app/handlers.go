package app

import (
	"encoding/json"
	"net/http"
)

type Customer struct {
	Name string `json:"name"`
}

func getAllCustomers(w http.ResponseWriter, r *http.Request) {
	customers := []Customer{
		{"Michael Street"},
		{"John Doe"},
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customers)

}
