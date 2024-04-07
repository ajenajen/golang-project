package handler

import (
	"bank/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type customerHandler struct {
	custService service.CustomerService
}

func NewCustomerHandler(custService service.CustomerService) customerHandler {
	return customerHandler{custService: custService}
}

func (h customerHandler) GetCustomers(w http.ResponseWriter, r *http.Request) { // w=writer
	customers, err := h.custService.GetCustomers()
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(customers)
}

func (h customerHandler) GetCustomer(w http.ResponseWriter, r *http.Request) {
	customerID, _ := strconv.Atoi(mux.Vars(r)["customerID"]) //strconv.Atoi ตัว convert str => int จะได้ return มา 2ค่า เป็น val, err
	// กรณีนี้เรามีการใส่ regX ดักไว้ก่อนจะเรียกแล้ว ว่าเป็นตัวเลขเท่านั้น ตอน convert เลยจะไม่เกิด error string

	customer, err := h.custService.GetCustomer(customerID)
	if err != nil {
		// appErr, ok := err.(errs.AppError) //เช็คว่าใช่ err ตัวข้างหน้า ใช่ type errs.AppError ไหม
		// if ok {
		// 	w.WriteHeader(appErr.Code)
		// 	fmt.Fprintln(w, appErr.Message)
		// 	return
		// }
		handleError(w, err)

		return
	}

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(customer)
}
