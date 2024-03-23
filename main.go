package main

type Person struct {
	Name string
	Age  int
}

// Extension method เป็นการประกาศ func ใน type
func (p Person) Hello() string {
	return "Hello " + p.Name
}

func main() {
	x := Person{Name: "Jane", Age: 18}
	println(x.Hello())
}
