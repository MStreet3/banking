package app

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mstreet3/banking/dto"
	"github.com/mstreet3/banking/service"
)

type TransactionHandlers struct {
	service service.TransactionService
}

func (th TransactionHandlers) newTransaction(w http.ResponseWriter, r *http.Request) {
	/* fetch URL params */
	vars := mux.Vars(r)
	id := vars["account_id"]

	var req dto.NewTransactionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	req.AccountId = id
	resp, appErr := th.service.NewTransaction(req)
	if appErr != nil {
		writeResponse(w, appErr.Code, appErr.AsMessage())
		return
	}
	writeResponse(w, http.StatusCreated, resp)
}
