package main

import (
	"bank/handler"
	"bank/logs"
	"bank/repository"
	"bank/service"
	"fmt"
	"net/http"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

func main() {
	initTimeZone()
	initConfig()

	db := initDatabase()
	customerRepositoryDB := repository.NewCustomerRepositoryDB(db)
	// customerRepositoryMock := repository.NewCustomerRepositoryMock() //ลองสลับไปใช้ data mockup ดูได้
	// _ = customerRepositoryMock
	customerService := service.NewCustomerService(customerRepositoryDB)
	customerHandler := handler.NewCustomerHandler(customerService)

	accountRepositoryDB := repository.NewAccountRepositoryDB(db)
	accountService := service.NewAccountService(accountRepositoryDB)
	accountHandler := handler.NewAccountHandler(accountService)

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

	router.HandleFunc("/customers/{customerID:[0-9]+}/accounts", accountHandler.NewAccount).Methods(http.MethodPost)
	router.HandleFunc("/customers/{customerID:[0-9]+}/accounts", accountHandler.GetAccounts).Methods(http.MethodGet)

	// log.Printf("Banking service started at port %v", viper.GetInt("app.port"))
	// logs.Log.Info("Banking service started at port " + viper.GetString("app.port"))
	// package.Log ที่ var ไว้.[info,debug,error] โดยตัว Info รับแต่ strin เลยต้องแปลงเป็น string หมด
	logs.Info("Banking service started at port " + viper.GetString("app.port"))

	http.ListenAndServe(fmt.Sprintf(":%v", viper.GetInt("app.port")), router)

}

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) //ใช้กรณี replace env ตอน start
	//> APP_PORT=3001 go run .

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func initTimeZone() {
	ict, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		panic(err)
	}

	time.Local = ict
}

func initDatabase() *sqlx.DB {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v",
		viper.GetString("db.username"),
		viper.GetString("db.password"),
		viper.GetString("db.host"),
		viper.GetInt("db.port"),
		viper.GetString("db.database"),
	)

	db, err := sqlx.Open(viper.GetString("db.driver"), dsn)
	if err != nil {
		panic(err)
	}

	db.SetConnMaxLifetime(3 * time.Minute)
	db.SetMaxOpenConns(10) // เปิดพร้อมกันได้กี่อัน
	db.SetMaxIdleConns(10)

	return db
}
