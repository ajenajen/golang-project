package services_test

import (
	"fmt"
	"gotest/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChekGrade(t *testing.T) {

	type testCase struct {
		name     string
		score    int
		expected string
	}

	cases := []testCase{
		{name: "A", score: 80, expected: "A"},
		{name: "B", score: 70, expected: "B"},
		{name: "C", score: 60, expected: "C"},
		{name: "D", score: 50, expected: "D"},
		{name: "F", score: 0, expected: "F"},
	}

	for _, c := range cases {
		// ใช้ t.Run ทำ sub test จะได้ run ต่อกันได้เลย
		t.Run(c.name, func(t *testing.T) {
			grade := services.CheckGrade(c.score) // กด cmd+. จะสร้าง func ให้
			// expected := c.expected

			assert.Equal(t, c.expected, grade)
			// if grade != expected {
			// 	t.Errorf("got %v expected %v", grade, expected)
			// }
		})
	}

}

// > go test gotest/services -bench=.
// ดู cpu mem ที่ใช้ด้วย > go test gotest/services -bench=. -benchmem
func BenchmarkCheckGrade(b *testing.B) {

	//b.N จำนวนรอบของการทดสอบ
	for i := 0; i < b.N; i++ {
		services.CheckGrade(80)
	}

}

// > go install golang.org/x/tools/cmd/godoc
// > godoc -http=:8000
// ใช้เพื่อเป็นตัวอย่าง ในการเรียกใช้ / มันจะไปขึ้นที่ document
func ExampleCheckGrade() {
	grade := services.CheckGrade(80)
	fmt.Println(grade)
	// Output: A
}

// http://localhost:8000/pkg/gotest/services/
