package series

import (
	"fmt"
	"time"

	"github.com/ptiger10/pd/options"
)

func ExampleSeries_empty_series() {
	s := newEmptySeries()
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

func ExampleNew_string_named() {
	s := MustNew([]string{"foo", "bar", "", "baz"}, Config{Name: "foobar"})
	fmt.Println(s)
	// Output:
	// 0    foo
	// 1    bar
	// 2    NaN
	// 3    baz
	// datatype: string
	// name: foobar
}

func ExampleNew_named_later() {
	s := MustNew(
		[]string{"foo", "bar", "", "baz"},
	)
	s.Rename("foobar")
	fmt.Println(s)

	// Output:
	// 0    foo
	// 1    bar
	// 2    NaN
	// 3    baz
	// datatype: string
	// name: foobar
}

func ExampleNew_string_singleIdx() {
	s := MustNew([]string{"foo", "bar", "", "baz"}, Config{Index: []int{100, 101, 102, 103}})
	fmt.Println(s)

	// Output:
	// 100    foo
	// 101    bar
	// 102    NaN
	// 103    baz
	// datatype: string
}

func ExampleNew_string_multiIdx() {
	s := MustNew(
		[]string{"foo", "bar", "", "baz"},
		Config{MultiIndex: []interface{}{[]int{0, 1, 2, 3}, []int{100, 101, 102, 103}}})
	fmt.Println(s)

	// Output:
	// 0 100    foo
	// 1 101    bar
	// 2 102    NaN
	// 3 103    baz
	// datatype: string
}

func ExampleNew_string_multiIdx_named_sequential_repeating() {
	s := MustNew(
		[]string{"foo", "bar", "", "baz"},
		Config{
			Name:            "foobar",
			MultiIndex:      []interface{}{[]int{0, 0, 1, 1}, []int{10000, 10100, 10200, 10300}},
			MultiIndexNames: []string{"id", "code"},
		},
	)
	fmt.Println(s)

	// Output:
	// id  code
	//  0 10000    foo
	//    10100    bar
	//  1 10200    NaN
	//    10300    baz
	// datatype: string
	// name: foobar
}

func ExampleNew_nonsequential_repeating() {
	s := MustNew(
		[]string{"foo", "bar", "", "baz", "qux", "quux"},
		Config{
			Name:            "foobar",
			MultiIndex:      []interface{}{[]int{0, 0, 1, 1, 0, 0}, []int{10000, 10100, 10200, 10300, 10400, 10500}},
			MultiIndexNames: []string{"id", "code"},
		},
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
	// name: foobar
}

func ExampleNew_partially_named_indexes() {
	s := MustNew(
		[]string{"foo", "bar"},
		Config{
			MultiIndex: []interface{}{
				[]int{0, 1},
				[]int{10000, 10100}},
			MultiIndexNames: []string{"", "code"},
			Name:            "foobar"},
	)
	fmt.Println(s)

	// Output:
	//    code
	// 0 10000    foo
	// 1 10100    bar
	// datatype: string
	// name: foobar
}

func ExampleNew_datetime_single() {
	s := MustNew([]time.Time{time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)})
	fmt.Println(s)

	// Output:
	// 0    1/1/2019T00:00:00
	// datatype: dateTime
}

func ExampleNew_config_nameOnly() {
	s := MustNew([]string{"foo", "bar"}, Config{Name: "baz"})
	fmt.Println(s)
	// Output:
	// 0    foo
	// 1    bar
	// datatype: string
	// name: baz
}

func ExampleNew_config_indexName() {
	s := MustNew([]string{"foo", "bar"}, Config{IndexName: "baz"})
	fmt.Println(s)
	// Output:
	// baz
	//   0    foo
	//   1    bar
	// datatype: string
}

func ExampleNew_config_datatype() {
	s := MustNew([]interface{}{"1", "foo"}, Config{DataType: options.Float64})
	fmt.Println(s)
	// Output:
	// 0    1
	// 1    NaN
	// datatype: float64
}

func ExampleElement_valid_printer() {
	s := MustNew("foo")
	fmt.Println(s.Element(0))

	// Output:
	//      Value: foo
	//       Null: false
	//     Labels: [0]
	// LabelTypes: [int64]
}

func ExampleElement_null_printer() {
	s := MustNew("")
	fmt.Println(s.Element(0))

	// Output:
	//      Value: NaN
	//       Null: true
	//     Labels: [0]
	// LabelTypes: [int64]
}