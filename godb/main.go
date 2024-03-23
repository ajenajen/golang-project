package main

import (
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Cover struct {
	Id   int
	Name string
}

var db *sqlx.DB // Use SQLX

func main() {
	var err error
	// db, err = sqlx.Open("sqlserver", "sqlserver://sa:P@ssw0rd@13.76.163.73?database=techcoach")
	db, err = sqlx.Open("mysql", "root:12345678@tcp(localhost:3306)/mockup_project")
	if err != nil {
		panic(err) //หยุดโปรแกรม
	}

	//Add
	// newCover := Cover{9, "cover-jane"}
	// err = AddCover(newCover)
	// if err != nil {
	// 	panic(err)
	// }

	//Update
	// cover := Cover{8, "cover-janejane"}
	// err = UpdateCover(cover)
	// if err != nil {
	// 	panic(err)
	// }

	//Delete
	// err = DeleteCover(8)
	// if err != nil {
	// 	panic(err)
	// }

	covers, err := GetCovers()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, cover := range covers {
		fmt.Println(cover)
	}

	cover, err := GetCover(7)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(cover)
}

func GetCovers() ([]Cover, error) {
	query := "select id, name from cover"
	covers := []Cover{}
	err := db.Select(&covers, query)
	if err != nil {
		return nil, err
	}
	// Select จะจับข้อมูลให้หมดเลย ไม่ต้องมา loop เอง ไม่ต้อง Scan อีก
	return covers, nil
}

func GetCover(id int) (*Cover, error) {
	query := "select id, name from cover where id=?"
	cover := Cover{}
	err := db.Get(&cover, query, id)
	if err != nil {
		return nil, err
	}
	// Get จะจับข้อมูลให้ ไม่ต้อง Scan อีก
	return &cover, nil
}

func AddCover(cover Cover) error {
	tx, err := db.Begin() // เหมือน memo ไว้ก่อน เผื่อพังตอนก่อน transaction จบ จะได้ไม่โดนลบ หรือ insert ไปก่อน ใช้ได้บาง driver
	if err != nil {       // ต้อง handle error ตอนเริ่มด้วย
		return err
	}

	query := "insert into cover (id, name) values (?, ?)"
	result, err := tx.Exec(query, cover.Id, cover.Name) //ใช้ tx แทน
	if err != nil {
		return err
	}

	affect, err := result.RowsAffected()
	if err != nil {
		tx.Rollback() // เผื่อมีอะไรไม่ผ่าน
		return err
	}
	if affect <= 0 {
		return errors.New("cannot insert")
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func UpdateCover(cover Cover) error {
	query := "update cover set name=? where id=?"
	result, err := db.Exec(query, cover.Name, cover.Id) //สลับตำแหน่ง argument ด้วย
	if err != nil {
		return err
	}

	affect, err := result.RowsAffected() // ค่า affect MySQL จะ return มาเป็นมากกว่า 0 ถ้าสำเร็จ
	if err != nil {
		return err
	}
	if affect <= 0 {
		return errors.New("cannot update")
	}
	return nil
}

func DeleteCover(id int) error {
	query := "delete from cover where id=?"
	result, err := db.Exec(query, id) //สลับตำแหน่ง argument ด้วย
	if err != nil {
		return err
	}

	affect, err := result.RowsAffected() // ค่า affect MySQL จะ return มาเป็นมากกว่า 0 ถ้าสำเร็จ
	if err != nil {
		return err
	}
	if affect <= 0 {
		return errors.New("cannot delete")
	}
	return nil
}
