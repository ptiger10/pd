package series

import (
	"fmt"

	"github.com/ptiger10/pd/opt"
)

func ExampleSeries_empty() {
	s := Series{}
	fmt.Println(s)

	// Output:
	// Series{}
}

func ExampleSeries_string_defaultIdx() {
	s, _ := New([]string{"Joe", "Jamy", "", "Jenny"})
	fmt.Println(s)

	// Output:
	// 0    Joe
	// 1    Jamy
	// 2    NaN
	// 3    Jenny
	// kind: string
}

func ExampleSeries_string_customIdx() {
	s, _ := New([]string{"Joe", "Jamy", "", "Jenny"}, Idx([]int{100, 101, 102, 103}))
	fmt.Println(s)

	// Output:
	// 100    Joe
	// 101    Jamy
	// 102    NaN
	// 103    Jenny
	// kind: string
}

func ExampleSeries_string_customIndex2() {
	s, _ := New(
		[]string{"Joe", "Jamy", "", "Jenny"},
		Idx([]int{0, 10, 11, 12}))
	fmt.Println(s)

	// Output:
	//  0    Joe
	// 10    Jamy
	// 11    NaN
	// 12    Jenny
	// kind: string
}

func ExampleSeries_string_multiIdx() {
	s, _ := New([]string{"Joe", "Jamy", "", "Jenny"}, Idx([]int{0, 1, 2, 3}), Idx([]int{100, 101, 102, 103}))
	fmt.Println(s)

	// Output:
	// 0 100    Joe
	// 1 101    Jamy
	// 2 102    NaN
	// 3 103    Jenny
	// kind: string
}

func ExampleSeries_string_namedindex_1() {
	s, _ := New(
		[]string{"Joe", "Jamy", "", "Jenny"}, opt.Name("student"),
		Idx([]int{0, 0, 1, 1}, opt.Name("id")),
		Idx([]int{10000, 10100, 10200, 10300}, opt.Name("code")),
	)
	fmt.Println(s)

	// Output:
	// id  code
	//  0 10000    Joe
	//    10100    Jamy
	//  1 10200    NaN
	//    10300    Jenny
	// kind: string
	// name: student
}

func ExampleSeries_string_namedindex_2() {
	s, _ := New(
		[]string{"Joe", "Jamy", "", "Jenny"}, opt.Name("student"),
		Idx([]int{100, 100, 101, 101}, opt.Name("id")),
		Idx([]int{10, 20, 30, 40}, opt.Name("code")),
	)
	fmt.Println(s)

	// Output:
	//  id code
	// 100   10    Joe
	//       20    Jamy
	// 101   30    NaN
	//       40    Jenny
	// kind: string
	// name: student
}

func ExampleSeries_string_named_index_nonsequential() {
	s, _ := New(
		[]string{"Joe", "Jamy", "", "Jenny", "Jeremiah", "Jemma"}, opt.Name("student"),
		Idx([]int{0, 0, 1, 1, 0, 0}, opt.Name("id")),
		Idx([]int{10000, 10100, 10200, 10300, 10400, 10500}, opt.Name("code")),
	)
	fmt.Println(s)

	// Output:
	// id  code
	//  0 10000    Joe
	//    10100    Jamy
	//  1 10200    NaN
	//    10300    Jenny
	//  0 10400    Jeremiah
	//    10500    Jemma
	// kind: string
	// name: student
}

func ExampleSeries_string_partial_named_Idx() {
	s, _ := New(
		[]string{"Joe", "Jamy", "", "Jenny"}, opt.Name("student"),
		Idx([]int{0, 0, 1, 1}),
		Idx([]int{10000, 10100, 10200, 10300}, opt.Name("code")),
	)
	fmt.Println(s)

	// Output:
	//    code
	// 0 10000    Joe
	//   10100    Jamy
	// 1 10200    NaN
	//   10300    Jenny
	// kind: string
	// name: student
}

// func TestPrinter(t *testing.T) {
// 	s, _ := New(
// 		[]string{"Joe", "Jamy", "", "Jenny", "Jeremiah", "Jemma"}, opt.Name("student"),
// 		Idx([]int{0, 0, 1, 1, 0, 0}, opt.Name("code")),
// 		Idx([]int{10000, 10100, 10200, 10300, 10400, 10500}),
// 	)
// 	fmt.Println(s.print())
// }

func ExampleElem() {
	s, _ := New("foo")
	fmt.Println(s.Element(0))

	// Output:
	//      Value: foo
	//       Null: false
	//     Labels: [0]
	// LabelKinds: [int64]
}
