package series

import (
	"fmt"
)

func ExampleSeries_empty() {
	s := Series{}
	fmt.Println(s)

	// Output:
	// Series{}
}

func ExampleNew_string() {
	s, _ := New([]string{"foo", "bar", "", "baz"})
	fmt.Println(s)

	// Output:
	// 0    foo
	// 1    bar
	// 2    NaN
	// 3    baz
	// datatype: string
}

// func ExampleNewWithConfig_string_named() {
// 	s, _ := NewWithConfig([]string{"foo", "bar", "", "baz"}, Config{Name: "foobar"})
// 	fmt.Println(s)

// 	// Output:
// 	// 0    foo
// 	// 1    bar
// 	// 2    NaN
// 	// 3    baz
// 	// datatype: string
// 	// name: foobar
// }

// func ExampleNew_named_later() {
// 	s, _ := New(
// 		[]string{"foo", "bar", "", "baz"},
// 	)
// 	s.Name = "foobar"
// 	fmt.Println(s)

// 	// Output:
// 	// 0    foo
// 	// 1    bar
// 	// 2    NaN
// 	// 3    baz
// 	// datatype: string
// 	// name: foobar
// }

// func ExampleNewWithConfig_string_singleIdx() {
// 	s, _ := NewWithConfig([]string{"foo", "bar", "", "baz"}, Config{Index: []int{100, 101, 102, 103}})
// 	fmt.Println(s)

// 	// Output:
// 	// 100    foo
// 	// 101    bar
// 	// 102    NaN
// 	// 103    baz
// 	// datatype: string
// }

// func ExampleNewWithConfig_string_multiIdx() {
// 	s, _ := NewWithConfig(
// 		[]string{"foo", "bar", "", "baz"},
// 		Config{MultiIndex: []interface{}{[]int{0, 1, 2, 3}, []int{100, 101, 102, 103}}})
// 	fmt.Println(s)

// 	// Output:
// 	// 0 100    foo
// 	// 1 101    bar
// 	// 2 102    NaN
// 	// 3 103    baz
// 	// datatype: string
// }

// func ExampleNewWithConfig_string_multiIdx_named_sequential_repeating() {
// 	s, _ := NewWithConfig(
// 		[]string{"foo", "bar", "", "baz"},
// 		Config{
// 			Name:            "foobar",
// 			MultiIndex:      []interface{}{[]int{0, 0, 1, 1}, []int{10000, 10100, 10200, 10300}},
// 			MultiIndexNames: []string{"id", "code"},
// 		},
// 	)
// 	fmt.Println(s)

// 	// Output:
// 	// id  code
// 	//  0 10000    foo
// 	//    10100    bar
// 	//  1 10200    NaN
// 	//    10300    baz
// 	// datatype: string
// 	// name: foobar
// }

// func ExampleNewWithConfig_nonsequential_repeating() {
// 	s, _ := NewWithConfig(
// 		[]string{"foo", "bar", "", "baz", "qux", "quux"},
// 		Config{
// 			Name:            "foobar",
// 			MultiIndex:      []interface{}{[]int{0, 0, 1, 1, 0, 0}, []int{10000, 10100, 10200, 10300, 10400, 10500}},
// 			MultiIndexNames: []string{"id", "code"},
// 		},
// 	)
// 	fmt.Println(s)

// 	// Output:
// 	// id  code
// 	//  0 10000    foo
// 	//    10100    bar
// 	//  1 10200    NaN
// 	//    10300    baz
// 	//  0 10400    qux
// 	//    10500    quux
// 	// datatype: string
// 	// name: foobar
// }

// func ExampleNewWithConfig_partially_named_indexes() {
// 	s, _ := NewWithConfig(
// 		[]string{"foo", "bar"},
// 		Config{
// 			MultiIndex: []interface{}{
// 				[]int{0, 1},
// 				[]int{10000, 10100}},
// 			MultiIndexNames: []string{"", "code"},
// 			Name:            "foobar"},
// 	)
// 	fmt.Println(s)

// 	// Output:
// 	//    code
// 	// 0 10000    foo
// 	// 1 10100    bar
// 	// datatype: string
// 	// name: foobar
// }

// func ExampleNew_datetime() {
// 	s, _ := New([]time.Time{time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)})
// 	fmt.Println(s)

// 	// Output:
// 	// 0    1/1/2019T00:00:00
// 	// datatype: dateTime
// }

// func ExampleElem_valid() {
// 	s, _ := New("foo")
// 	fmt.Println(s.Element(0))

// 	// Output:
// 	//      Value: foo
// 	//       Null: false
// 	//     Labels: [0]
// 	// LabelTypes: [int64]
// }

// func ExampleElem_null() {
// 	s, _ := New("")
// 	fmt.Println(s.Element(0))

// 	// Output:
// 	//      Value: NaN
// 	//       Null: true
// 	//     Labels: [0]
// 	// LabelTypes: [int64]
// }
