package series

import (
	"fmt"
	"time"
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
	// datatype: string
}

func ExampleSeries_string_customIdx() {
	s, _ := New([]string{"Joe", "Jamy", "", "Jenny"}, Idx([]int{100, 101, 102, 103}))
	fmt.Println(s)

	// Output:
	// 100    Joe
	// 101    Jamy
	// 102    NaN
	// 103    Jenny
	// datatype: string
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
	// datatype: string
}

func ExampleSeries_string_multiIdx() {
	s, _ := New([]string{"Joe", "Jamy", "", "Jenny"}, Idx([]int{0, 1, 2, 3}), Idx([]int{100, 101, 102, 103}))
	fmt.Println(s)

	// Output:
	// 0 100    Joe
	// 1 101    Jamy
	// 2 102    NaN
	// 3 103    Jenny
	// datatype: string
}

func ExampleNewWithConfig_string() {
	s, _ := NewWithConfig(
		Config{Name: "student"}, []string{"foo", "bar", "", "baz"},
		IndexLevel{Labels: []int{0, 0, 1, 1}, Name: "id"},
		IndexLevel{Labels: []int{10000, 10100, 10200, 10300}, Name: "code"},
	)
	fmt.Println(s)

	// Output:
	// id  code
	//  0 10000    foo
	//    10100    bar
	//  1 10200    NaN
	//    10300    baz
	// datatype: string
	// name: student
}

func ExampleNewWithConfig_nonsequential_repeating() {
	s, _ := NewWithConfig(
		Config{Name: "student"}, []string{"foo", "bar", "", "baz", "qux", "quux"},
		IndexLevel{Labels: []int{0, 0, 1, 1, 0, 0}, Name: "id"},
		IndexLevel{Labels: []int{10000, 10100, 10200, 10300, 10400, 10500}, Name: "code"},
	)
	fmt.Println(s)

	// Output:
	// id  code
	//  0 10000    foo
	//    10100    bar
	//  1 10200    NaN
	//    10300    baz
	//  0 10400    qux
	//    10500    quux
	// datatype: string
	// name: student
}

func ExampleNewWithConfig_partially_named_indexes() {
	s, _ := NewWithConfig(
		Config{Name: "student"}, []string{"foo", "bar", "", "baz", "qux", "quux"},
		IndexLevel{Labels: []int{0, 0, 1, 1, 0, 0}, Name: "id"},
		IndexLevel{Labels: []int{10000, 10100, 10200, 10300, 10400, 10500}, Name: "code"},
	)
	fmt.Println(s)

	// Output:
	//    code
	// 0 10000    Joe
	//   10100    Jamy
	// 1 10200    NaN
	//   10300    Jenny
	// datatype: string
	// name: student
}

func ExampleSeries_datetime() {
	s, _ := New([]time.Time{time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)})
	fmt.Println(s)

	// Output:
	// 0    1/1/2019T00:00:00
	// datatype: dateTime
}

func ExampleElem() {
	s, _ := New("foo")
	fmt.Println(s.Element(0))

	// Output:
	//      Value: foo
	//       Null: false
	//     Labels: [0]
	// LabelKinds: [int64]
}
