package app

import "net/http"

func Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/customers", getAllCustomers)
	http.ListenAndServe("localhost:8080", mux)
}
