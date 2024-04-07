package handler

import (
	"bank/errs"
	"bank/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type accountHandler struct {
	accService service.AccountService
}

func NewAccountHandler(accService service.AccountService) accountHandler {
	return accountHandler{accService: accService}
}

func (h accountHandler) NewAccount(w http.ResponseWriter, r *http.Request) {
	customerID, _ := strconv.Atoi(mux.Vars(r)["customerID"])

	if r.Header.Get("content-type") != "application/json" {
		// handle technical error
		handleError(w, errs.NewValidationError("request body incorrect format"))
		return
	}

	request := service.NewAccountRequest{}
	err := json.NewDecoder(r.Body).Decode(&request) // check decode body ว่าตรง type ที่กำหนดไหม
	if err != nil {
		handleError(w, errs.NewValidationError("request body incorrect format"))
		return
	}

	response, err := h.accService.NewAccount(customerID, request)
	if err != nil {
		handleError(w, err) // handle business error ซึ่่ง busi เขา handle มาแล้ว ก็เลยส่งต่อได้เลย
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h accountHandler) GetAccounts(w http.ResponseWriter, r *http.Request) {
	customerID, _ := strconv.Atoi(mux.Vars(r)["customerID"])

	responses, err := h.accService.GetAccounts(customerID)
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(responses) //ตรงนี้อาจจะมี error ต้องไป handle ต่อด้วย
}
