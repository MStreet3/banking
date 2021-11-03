package app

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Start() {
	router := mux.NewRouter()
	router.HandleFunc("/customers", getAllCustomers)
	http.ListenAndServe("localhost:8080", router)
}
