package series

import (
	"fmt"
)

func ExampleSeries_string_defaultIndex() {
	s, _ := New([]string{"Joe", "Jamy", "", "Jenny"})
	fmt.Println(s)

	// Output:
	// 0    Joe
	// 1    Jamy
	// 2
	// 3    Jenny
	// kind: string
}

func ExampleSeries_string_customIndex() {
	s, _ := New([]string{"Joe", "Jamy", "", "Jenny"}, Index([]int{100, 101, 102, 103}))
	fmt.Println(s)

	// Output:
	// 100    Joe
	// 101    Jamy
	// 102
	// 103    Jenny
	// kind: string
}

func ExampleSeries_string_multiIndex() {
	s, _ := New([]string{"Joe", "Jamy", "", "Jenny"}, Index([]int{0, 1, 2, 3}), Index([]int{100, 101, 102, 103}))
	fmt.Println(s)

	// Output:
	// 0 100    Joe
	// 1 101    Jamy
	// 2 102
	// 3 103    Jenny
	// kind: string
}

func ExampleSeries_string_namedindex_1() {
	s, _ := New(
		[]string{"Joe", "Jamy", "", "Jenny"}, Name("student"),
		Index([]int{0, 0, 1, 1}, Name("id")),
		Index([]int{10000, 10100, 10200, 10300}, Name("code")),
	)
	fmt.Println(s)

	// Output:
	// id  code
	//  0 10000    Joe
	//    10100    Jamy
	//  1 10200
	//    10300    Jenny
	// kind: string
	// name: student
}

func ExampleSeries_string_namedindex_2() {
	s, _ := New(
		[]string{"Joe", "Jamy", "", "Jenny"}, Name("student"),
		Index([]int{100, 100, 101, 101}, Name("id")),
		Index([]int{10, 20, 30, 40}, Name("code")),
	)
	fmt.Println(s)

	// Output:
	//  id code
	// 100   10    Joe
	//       20    Jamy
	// 101   30
	//       40    Jenny
	// kind: string
	// name: student
}

func ExampleSeries_string_named_index_nonsequential() {
	s, _ := New(
		[]string{"Joe", "Jamy", "", "Jenny", "Jeremiah", "Jemma"}, Name("student"),
		Index([]int{0, 0, 1, 1, 0, 0}, Name("id")),
		Index([]int{10000, 10100, 10200, 10300, 10400, 10500}, Name("code")),
	)
	fmt.Println(s)

	// Output:
	// id  code
	//  0 10000    Joe
	//    10100    Jamy
	//  1 10200
	//    10300    Jenny
	//  0 10400    Jeremiah
	//    10500    Jemma
	// kind: string
	// name: student
}

func ExampleSeries_string_partial_named_index() {
	s, _ := New(
		[]string{"Joe", "Jamy", "", "Jenny"}, Name("student"),
		Index([]int{0, 0, 1, 1}),
		Index([]int{10000, 10100, 10200, 10300}, Name("code")),
	)
	fmt.Println(s)

	// Output:
	//    code
	// 0 10000    Joe
	//   10100    Jamy
	// 1 10200
	//   10300    Jenny
	// kind: string
	// name: student
}

// func TestPrinter(t *testing.T) {
// 	s, _ := New(
// 		[]string{"Joe", "Jamy", "", "Jenny", "Jeremiah", "Jemma"}, Name("student"),
// 		Index([]int{0, 0, 1, 1, 0, 0}, Name("code")),
// 		Index([]int{10000, 10100, 10200, 10300, 10400, 10500}),
// 	)
// 	fmt.Println(s.print())
// }
