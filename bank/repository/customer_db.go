package repository

import "github.com/jmoiron/sqlx"

type customerRepositoryDB struct { //เป็นตัวเล็กเพราะไม่ให้ Exposed ออกไปที่อื่นได้ จาก main ให้เข้ามาถึงได้แค่ interface
	db *sqlx.DB
}

func NewCustomerRepositoryDB(db *sqlx.DB) customerRepositoryDB { // จาก main ให้มาเอาที่นี่แทน ใช้ในการ new instant ของ struct นี้ขึ้นมา
	return customerRepositoryDB{db: db}
}

func (r customerRepositoryDB) GetAll() ([]Customer, error) {
	customers := []Customer{}
	query := "select customer_id, name, date_of_birth, city, zipcode, status from customers"
	err := r.db.Select(&customers, query)
	if err != nil {
		return nil, err
	}

	return customers, nil
}

func (r customerRepositoryDB) GetById(id int) (*Customer, error) {
	customer := Customer{}
	query := "select customer_id, name, date_of_birth, city, zipcode, status from customers where customer_id=?"
	err := r.db.Get(&customer, query, id)
	if err != nil {
		return nil, err
	}

	return &customer, nil
}
