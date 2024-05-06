package main

import (
	"errors"
	"fmt"

	"github.com/stretchr/testify/mock"
)

func main() {
	// fmt.Println("Hello Test")

	c := CustomerRepositoryMock{}
	//=== mock
	c.On("GetCustomer", 1).Return("Jane", 18, nil) //เมื่อมีคนใช้ GetCustomer ส่ง 1 เข้ามา จะ return ...
	c.On("GetCustomer", 2).Return("", 0, errors.New("not found"))

	//=== ใช้งาน
	name, age, err := c.GetCustomer(1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(name, age)
}

// Example basic mock
type CustomerRepositoryMock struct {
	mock.Mock
}

// ถ้าเมื่อไรเป็น Mock ต้องส่ง reciever func เป็น pointer repo เพราะมันจะปลอมตัวเป็น mock struct
func (m *CustomerRepositoryMock) GetCustomer(id int) (name string, age int, err error) {
	args := m.Called(id) //เรียกค่าจาก mock / return เป็น slice
	return args.String(0), args.Int(1), args.Error(2)
}
