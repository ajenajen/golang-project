package main

import (
	"bank/handler"
	"bank/repository"
	"bank/service"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func main() {
	db, err := sqlx.Open("mysql", "root:12345678@tcp(localhost:3306)/banking")
	if err != nil {
		panic(err)
	}

	customerRepositoryDB := repository.NewCustomerRepositoryDB(db)
	customerRepositoryMock := repository.NewCustomerRepositoryMock() //ลองสลับไปใช้ data mockup ดูได้
	_ = customerRepositoryMock
	customerService := service.NewCustomerService(customerRepositoryDB)
	customerHandler := handler.NewCustomerHandler(customerService)

	// customers, err := customerRepository.GetAll()
	// if err != nil {
	// 	panic(err)
	// }
	// for _, customer := range customers {
	// 	fmt.Println(customer)
	// }

	// customer, err := customerRepository.GetById(2000)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(customer)

	/////////// Service ///////////
	// customers, err := customerService.GetCustomers()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(customers)

	// customer, err := customerService.GetCustomer(2000)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(customer)

	/////////// REST Handler ///////////
	router := mux.NewRouter()

	router.HandleFunc("/customers", customerHandler.GetCustomers).Methods(http.MethodGet) // test by > curl localhost:8000/customers
	// > curl localhost:8000/customers -i จะเห็น header ด้วย
	router.HandleFunc("/customers/{customerID:[0-9]+}", customerHandler.GetCustomer).Methods(http.MethodGet)

	http.ListenAndServe(":8000", router)

	/////////// Service ///////////
}
