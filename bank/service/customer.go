package service

type CustomerResponse struct {
	CustomerID int    `json:"customer_id"` //converting เป็น json
	Name       string `json:"name"`
	Status     int    `json:"status"`
}

type CustomerService interface {
	GetCustomers() ([]CustomerResponse, error)
	GetCustomer(int) (*CustomerResponse, error)
}
