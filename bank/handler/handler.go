package handler

import (
	"bank/errs"
	"fmt"
	"net/http"
)

func handleError(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case errs.AppError:
		w.WriteHeader(e.Code)
		fmt.Fprintln(w, e.Message)
	case error: //ถ้าไม่ใช่ AppError ก็ return ได้อยุ
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, e)
	}
}
