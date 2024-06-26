package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/go-sql-driver/mysql"
	jwtware "github.com/gofiber/jwt/v2"
)

var db *sqlx.DB

const jwtSecret = "Infinitas"

func main() {

	var err error
	db, err = sqlx.Open("mysql", "root:1234@tcp(localhost:3306)/learngolang")
	if err != nil {
		panic(err)
	}

	app := fiber.New()

	app.Use("/hello", jwtware.New(jwtware.Config{
		SigningMethod: "HS256",
		SigningKey:    []byte(jwtSecret),
		SuccessHandler: func(c *fiber.Ctx) error {
			return c.Next()
		},
		ErrorHandler: func(c *fiber.Ctx, e error) error {
			return fiber.ErrUnauthorized
		},
	}))

	app.Post("/signup", Signup)
	app.Post("/login", Login)
	app.Get("/hello", Hello)

	app.Listen(":8000")
}

func Signup(c *fiber.Ctx) error {
	request := SignupRequest{}
	err := c.BodyParser(&request)
	if err != nil {
		return err
	}

	if request.Username == "" || request.Password == "" {
		return fiber.ErrUnprocessableEntity
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), 10)
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	query := "insert user (username, password) values (?, ?)"
	result, err := db.Exec(query, request.Username, string(password))
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	user := User{
		Id:       int(id),
		Username: request.Username,
		Password: string(password),
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

// curl localhost:8000/signup -H content-type:application/json -d '{"username": "jane", "password": "jane"}' -i

func Login(c *fiber.Ctx) error {
	request := LoginRequest{}
	err := c.BodyParser(&request)
	if err != nil {
		return err
	}

	if request.Username == "" || request.Password == "" {
		return fiber.ErrUnprocessableEntity
	}

	user := User{}
	query := "select id, username, password from user where username=?"
	err = db.Get(&user, query, request.Username) // ส่ง user เข้าไปรับ
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Incorrect username or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Incorrect username or password")
	}

	claims := jwt.StandardClaims{
		Issuer:    strconv.Itoa(user.Id),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwtToken.SignedString([]byte(jwtSecret))
	if err != nil {
		return fiber.ErrInternalServerError
	}

	// return c.SendStatus(fiber.StatusOK)
	return c.JSON(fiber.Map{
		"jwtToken": token,
	})
}

// curl localhost:8000/login -H content-type:application/json -d '{"username": "jane", "password": "jane"}' -i

func Hello(c *fiber.Ctx) error {
	return c.SendString("Hello World")
}

// curl localhost:8000/hello -H "Authorization:Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTI2MDY1NDUsImlzcyI6IjEifQ.lfSaiFM4QFi4lncTuNYBL7H4AUTcTeETQjL0rqh-yvI"

type MyCustomClaims struct {
	jwt.StandardClaims
}

type User struct {
	Id       int    `db:"id" json:"id"`
	Username string `db:"username" json:"username"`
	Password string `db:"password" json:"password"`
}

type SignupRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Basic fiber
func Fiber() {
	app := fiber.New(fiber.Config{
		Prefork: true, //optionเพื่อspawnตัวเองออกมารับงานใน same port... reuse port
	})

	//Middleware: ทำงานก่อนและ หลัง ของแต่ละ path ได้
	// app.Use(func(c *fiber.Ctx) error {
	// 	c.Locals("name", "jane") //ประกาศตัวแปร ไปใช้กับทุก path
	// 	fmt.Println("before")
	// 	err := c.Next() //สั่งให้ทำคำสั่งใน path ต่อไป ถ้ามี error ให้ return error ต่อด้วย
	// 	fmt.Println("after")
	// 	return err
	// })
	// แบบกำหนด path
	app.Use("/hello", func(c *fiber.Ctx) error {
		fmt.Println("before only hello")
		err := c.Next()
		return err
	})

	//github.com/gofiber/fiber/v2/middleware/requestid
	app.Use(requestid.New())
	//test curl curl localhost:8000/hello -i
	//จะได้ X-Request-Id: 3d20fd91-df21-4061-9446-64ab019dd3ab ออกมาด้วยที่ header
	//config header เพิ่มเติมได้

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "*",
		AllowHeaders: "*",
	}))
	//test curl localhost:8000/hello -i
	//จะได้ Access-Control-Allow-Origin: * ออกมาด้วยที่ header

	//Logger
	// app.Use(logger.New(logger.Config{
	// 	TimeZone: "Asia/Bangkok",
	// }))

	//GET
	app.Get("/hello", func(c *fiber.Ctx) error {
		name := c.Locals("name") //ใช้ค่าตัวแปร local ที่ตั้งไว้จาก middleware
		return c.SendString(fmt.Sprintf("GET: Hello %v", name))
	})

	//POST
	app.Post("/hello", func(c *fiber.Ctx) error {
		return c.SendString("POST: Hello world")
	})

	//Parameter (Optional :param?)
	app.Get("/hello/:name/:surname?", func(c *fiber.Ctx) error {
		name := c.Params("name")
		surname := c.Params("surname")
		return c.SendString("GET: Hello world name: " + name + ", surname: " + surname)
	})

	//ParamsInt
	app.Get("/hello/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		// ถ้า id เป็น string จะไปเข้า case ก่อนหน้านี้
		if err != nil {
			return fiber.ErrBadRequest
		}
		return c.SendString(fmt.Sprintf("ID = %v", id))
	})

	//Query, curl "localhost:8000/query?name=jane&surname=suwa" -i
	app.Get("/query", func(c *fiber.Ctx) error {
		name := c.Query("name")
		surname := c.Query("surname")
		return c.SendString("name: " + name + ", surname: " + surname)
	})

	//Query 2 parser, curl "localhost:8000/query2?id=1&name=jane" -i
	app.Get("/query2", func(c *fiber.Ctx) error {
		person := Person{}
		c.QueryParser(&person)
		return c.JSON(person) //c.JSON แค่กำหนด struct มันจะแปลง Marshal ให้หมดเลย
	})

	//Wildcards, curl "localhost:8000/wildcards/hello/world" -i
	app.Get("/wildcards/*", func(c *fiber.Ctx) error {
		wildcard := c.Params("*")
		return c.SendString("wildcards: " + wildcard)
	})

	//Static File: เรียกไปที่ static ใน wwwroot
	// app.Static("/", "./wwwroot")
	//chain ต่อได้
	app.Static("/", "./wwwroot", fiber.Static{
		Index:         "index.html",
		CacheDuration: time.Second * 10, //cache ไว้ 10 วิ
	})

	//NewError
	app.Get("/error", func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusNotFound, "content not found")
	})

	//Group : for versioning api
	v1 := app.Group("/v1", func(c *fiber.Ctx) error { //เปรียบเสมือน middleware ของ group
		c.Set("Version", "v1") //Set ในนี้ คือ set Header, จะได้ Version: v1 มาใน header ด้วย
		return c.Next()
	})
	v1.Get("/hello", func(c *fiber.Ctx) error {
		return c.SendString("Hello group v1")
	})

	v2 := app.Group("/v2")
	v2.Get("/hello", func(c *fiber.Ctx) error {
		return c.SendString("Hello group v2")
	})

	//Mount: เป็นการแตกตัวออกจาก app ออกมาเลย ไม่ได้ใช้ของ app ที่ประกาศไว้ตอนแรก
	// เหมือนไว้จัดการ microservice ตัวนั้นตัวเดียวเลย
	userApp := fiber.New()
	userApp.Get("/login", func(c *fiber.Ctx) error {
		return c.SendString("User Login")
	})

	app.Mount("/user", userApp) //อะไรที่ตามหลัง /user จะมาใช้ userApp

	//Server
	app.Server().MaxConnsPerIP = 1 //กำหนด connection per ip
	app.Get("/server", func(c *fiber.Ctx) error {
		time.Sleep(time.Second * 10)
		return c.SendString("server")
	}) // ถ้าไปยิงซ้ำระหว่างรอ The number of connections from your ip exceeds MaxConnsPerIP%

	//Environment
	app.Get("/env", func(c *fiber.Ctx) error {
		//Json Map ค่าออกมาเลยก็ได้ ไม่ต้อง struct ก่อน
		return c.JSON(fiber.Map{
			"BaseURL":     c.BaseURL(),
			"Hostname":    c.Hostname(),
			"IP":          c.IP(),
			"IPs":         c.IPs(),         // IPs จะมีให้เห็นว่าผ่าน proxy อะไรมาหมดเลย
			"OriginalURL": c.OriginalURL(), // ได้ queryString มาด้วย
			"Path":        c.Path(),
			"Protocol":    c.Protocol(),
			"Subdomains":  c.Subdomains(),
		})
	})

	//Body
	app.Post("/body", func(c *fiber.Ctx) error {
		fmt.Printf("IsJson: %v \n", c.Is("json")) //เช็คที่ header
		fmt.Println(string(c.Body()))
		return nil
	})
	//curl "localhost:8000/body" -d 'hello'
	//curl "localhost:8000/body" -H content-type:"application/json" -d '{"name":"jane"}'

	//Body Parser: อ่านค่าแล้วแปลงเป็น struct ให้เลย
	app.Post("body-parser", func(c *fiber.Ctx) error {
		person := Person{}
		err := c.BodyParser(&person) //err มาตรวจว่าแปลงได้ไหม
		if err != nil {
			return err
		}

		fmt.Println(person)
		return nil
	})
	//curl "localhost:8000/body-parser" -H content-type:"application/json" -d '{"id": 1, "name":"jane"}'
	//{1 jane}
	app.Post("body-parser2", func(c *fiber.Ctx) error {
		// person := Person{}
		person := map[string]interface{}{} //interface{} like any
		err := c.BodyParser(&person)       //err มาตรวจว่าแปลงได้ไหม
		if err != nil {
			return err
		}

		fmt.Println(person)
		return nil
	})
	//curl "localhost:8000/body-parser2" -H content-type:"application/json" -d '{"id": 1, "name":"jane"}'
	//map[id:1 name:jane]
	//person := Person{} ถ้ากำหนด struct ที่แน่นอนไว้ จะส่งไป type ไหนก็ได้ แต่ต้องส่ง header type ไปให้ตรงเพื่อ parser เช็ค
	//curl "localhost:8000/body-parser2" -H content-type:"application/x-www-form-urlencoded" -d 'id=1&name=jane'
	//map[id:1 name:jane]

	app.Listen(":8000")
}

type Person struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
