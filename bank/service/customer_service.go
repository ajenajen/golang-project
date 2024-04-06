package service

import (
	"bank/errs"
	"bank/logs"
	"bank/repository"
	"database/sql"
)

type customerService struct {
	customerRepo repository.CustomerRepository
}

func NewCustomerService(customerRepo repository.CustomerRepository) customerService {
	return customerService{customerRepo: customerRepo}
}

func (s customerService) GetCustomers() ([]CustomerResponse, error) {
	customers, err := s.customerRepo.GetAll()
	if err != nil {
		// log.Println(err)
		logs.Error(err)
		return nil, errs.NewUnexpectedError()
	}

	custResponses := []CustomerResponse{}
	for _, customer := range customers {
		custResponse := CustomerResponse{
			CustomerID: customer.CustomerID,
			Name:       customer.Name,
			Status:     customer.Status,
		}
		custResponses = append(custResponses, custResponse)
	}

	return custResponses, nil
}

func (s customerService) GetCustomer(id int) (*CustomerResponse, error) {
	customer, err := s.customerRepo.GetById(id)
	if err != nil {

		if err == sql.ErrNoRows {
			// return nil, errors.New("customer not found")
			// return nil, errs.AppError{
			// 	Code:    http.StatusNotFound,
			// 	Message: "customer not found",
			// }
			return nil, errs.NewNotfoundError("customer not found")
		}

		// log.Println(err)
		logs.Error(err)
		// return nil, err
		// return nil, errs.AppError{
		// 	Code:    http.StatusInternalServerError,
		// 	Message: "unexpected error",
		// }
		return nil, errs.NewUnexpectedError()
	}

	custResponse := CustomerResponse{
		CustomerID: customer.CustomerID,
		Name:       customer.Name,
		Status:     customer.Status,
	}

	return &custResponse, nil
}
