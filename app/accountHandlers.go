package app

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mstreet3/banking/dto"
	"github.com/mstreet3/banking/service"
	"github.com/mstreet3/banking/utils"
)

type AccountHandlers struct {
	service service.AccountService
}

func (ah AccountHandlers) newAccount(w http.ResponseWriter, r *http.Request) {
	/* fetch URL params */
	vars := mux.Vars(r)
	id := vars["customer_id"]

	var req dto.NewAccountRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.WriteResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	req.CustomerId = id
	resp, appErr := ah.service.NewAccount(req)
	if appErr != nil {
		utils.WriteResponse(w, appErr.Code, appErr.AsMessage())
		return
	}
	utils.WriteResponse(w, http.StatusCreated, resp)

}

func (ah AccountHandlers) newTransaction(w http.ResponseWriter, r *http.Request) {
	/* fetch URL params */
	vars := mux.Vars(r)
	id := vars["account_id"]

	var req dto.NewTransactionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.WriteResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	req.AccountId = id
	resp, appErr := ah.service.MakeTransaction(req)
	if appErr != nil {
		utils.WriteResponse(w, appErr.Code, appErr.AsMessage())
		return
	}
	utils.WriteResponse(w, http.StatusCreated, resp)
}
