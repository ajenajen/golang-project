package repository

import "errors"

type customerRepositoryMock struct {
	customers []Customer
}

func NewCustomerRepositoryMock() customerRepositoryMock {
	customers := []Customer{
		{CustomerID: 1000, Name: "John Doe", DateOfBirth: "1990-05-15", City: "New York", ZipCode: "10001", Status: 1},
		{CustomerID: 1001, Name: "Jane Smith", DateOfBirth: "1985-09-25", City: "Los Angeles", ZipCode: "90001", Status: 1},
	}

	return customerRepositoryMock{customers: customers}
}

func (r customerRepositoryMock) GetAll() ([]Customer, error) {
	return r.customers, nil
}

func (r customerRepositoryMock) GetById(id int) (*Customer, error) {
	for _, cutomer := range r.customers {
		if cutomer.CustomerID == id {
			return &cutomer, nil
		}
	}

	return nil, errors.New("customer not found")
}
