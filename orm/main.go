package main

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type SqlLogger struct {
	logger.Interface // ทำให้ SqlLogger มัน conform ตาม type Interface ของ gorm เพื่อจะเอา Trace ไปใช้ log sql statement
}

func (l SqlLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, _ := fc()
	fmt.Printf("%v\n===============================================\n", sql)
}

var db *gorm.DB

func main() {
	dsn := "root:1234@tcp(localhost:3306)/learngolang?parseTime=true" //https://pkg.go.dev/github.com/go-sql-driver/mysql
	dial := mysql.Open(dsn)

	var err error
	db, err = gorm.Open(dial, &gorm.Config{
		Logger: &SqlLogger{},
		// DryRun: true, // จะ log sql script ออกมาดูก่อน ยังไม่สร้างจริง
	})
	//ถ้า gorm ทำงาน จะส่ง log ผ่าน struct SqlLogger ที่เราสร้างไว้ด้วย
	if err != nil {
		panic(err)
	}

	// db.Migrator().CreateTable(Test{}) // create table ให้เลย
	// db.AutoMigrate(Gender{}, Test{}) // gorm จะเช็คก่อนว่ามี table เดิมอยู่ไหม ก่อนสร้าง จะไม่ error
	// CreateGender("xxxx")
	// GetGenders()
	// GetGender(1)
	// GetGenderByName("Male")

	// UpdateGender(4, "yyyy")
	// `
	// SELECT * FROM `genders` WHERE `genders`.`id` = 4 ORDER BY `genders`.`id` LIMIT 1
	// ===============================================
	// UPDATE `genders` SET `name`='yyyy' WHERE `id` = 4
	// ===============================================
	// SELECT * FROM `genders` WHERE `genders`.`id` = 4 ORDER BY `genders`.`id` LIMIT 1
	// `

	// UpdateGender2(4, "zzzz")
	// `
	// UPDATE `genders` SET `name`='zzzz' WHERE id=4
	// ===============================================
	// SELECT * FROM `genders` WHERE `genders`.`id` = 4 ORDER BY `genders`.`id` LIMIT 1
	// `
	//UpdateGender2(4, "") //ส่งแบบนี้เข้าไปที่ท่านี้ มันไม่ทำงานนะ bugเลย

	// DeleteGender(5)
	// `
	// DELETE FROM `genders` WHERE `genders`.`id` = 5
	// ===============================================
	// Deleted
	// SELECT * FROM `genders` WHERE `genders`.`id` = 5 ORDER BY `genders`.`id` LIMIT 1
	// ===============================================
	// record not found
	// `

	// CreateTest(0, "Test 1")
	// `
	// UPDATE `MyTest` SET `deleted_at`='2024-04-08 07:41:02.732' WHERE `MyTest`.`id` = 3 AND `MyTest`.`deleted_at` IS NULL
	// `
	// CreateTest(0, "Test 2")
	// CreateTest(0, "Test 3")

	// DeleteTest(3)
	// ` มันเจอ deleteAt เลยไปสั่ง update เป็น soft delete ให้
	// UPDATE `MyTest` SET `deleted_at`='2024-04-08 07:41:02.732' WHERE `MyTest`.`id` = 3 AND `MyTest`.`deleted_at` IS NULL
	// `

	// DeletePermanentTest(3) //DELETE FROM `MyTest` WHERE `MyTest`.`id` = 3

	// GetTests()
	//SELECT * FROM `MyTest` WHERE `MyTest`.`deleted_at` IS NULL
}

func GetGenders() {
	genders := []Gender{}
	tx := db.Order("id").Find(&genders) //SELECT * FROM `genders`
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println(genders)
}

func GetGender(id uint) {
	gender := Gender{}
	tx := db.First(&gender, id) //SELECT * FROM `genders` WHERE `genders`.`id` = 1 ORDER BY `genders`.`id` LIMIT 1
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println(gender)
}

func GetGenderByName(name string) {
	gender := Gender{}
	// tx := db.Find(&gender, "name=?", name)
	tx := db.Where("name=?", name).Find(&gender)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println(gender)
}

func UpdateGender(id uint, name string) {
	gender := Gender{}
	tx := db.First(&gender, id) // check in db ก่อน
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}

	gender.Name = name
	tx = db.Save(&gender)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	GetGender(id)
}

func UpdateGender2(id uint, name string) {
	gender := Gender{Name: name}                                // update แบบนี้ ค่าเดิมต้องไม่ใช่ 0 ด้วย
	tx := db.Model(&Gender{}).Where("id=?", id).Updates(gender) //ต้องระบุ struct ที่จะอัพเดทเข้าไปที่ model
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}

	GetGender(id)
}

func DeleteGender(id uint) {
	tx := db.Delete(&Gender{}, id)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}

	fmt.Println("Deleted")
	GetGender(id)
}

func CreateGender(name string) {
	gender := Gender{
		Name: name,
	}
	tx := db.Create(&gender)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println(gender)
}

type Gender struct {
	ID   uint
	Name string `gorm:"unique; size:10"`
}

type Test struct {
	gorm.Model // มี ID, CreateAt, UpdateAt, DeleteAt มาให้
	Code       uint
	Name       string
	//CREATE TABLE `tests` (`id` bigint unsigned AUTO_INCREMENT,`created_at` datetime(3) NULL,`updated_at` datetime(3) NULL,`deleted_at` datetime(3) NULL,`code` bigint unsigned,`name` longtext,PRIMARY KEY (`id`),INDEX `idx_tests_deleted_at` (`deleted_at`))
}

func (t Test) TableName() string {
	return "MyTest"
	//CREATE TABLE `MyTest`
}

type TestGorm struct {
	ID   uint
	Code uint   `gorm:"primaryKey;comment: This is Code"`
	Name string `gorm:"column:myname;type:varchar(50);unique;default:Hello;not null"`
	//string type default 'longtext'
	//ใช้ type:varchar(50) or size:50 ก็ได้
	//CREATE TABLE `test_gorms` (`id` bigint unsigned,`code` bigint unsigned AUTO_INCREMENT COMMENT ' This is Code',`myname` varchar(50) NOT NULL DEFAULT 'Hello',PRIMARY KEY (`code`),CONSTRAINT `uni_test_gorms_myname` UNIQUE (`myname`))
}

func CreateTest(code uint, name string) {
	test := Test{
		Code: code,
		Name: name,
	}
	db.Create(&test)
}

func GetTests() {
	tests := []Test{}
	db.Find(&tests)

	for _, test := range tests {
		fmt.Printf("%v|%v \n", test.ID, test.Name)
	}
}

func DeleteTest(id uint) {
	db.Delete(&Test{}, id)
}

func DeletePermanentTest(id uint) {
	db.Unscoped().Delete(&Test{}, id)
}
