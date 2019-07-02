package series

import (
	"fmt"
	"math"
	"time"

	"github.com/ptiger10/pd/options"
)

// [START constructor examples]
func ExampleNew_string() {
	s, _ := New([]string{"foo", "bar", "", "baz"})
	fmt.Println(s)
	// Output:
	// 0    foo
	// 1    bar
	// 2    NaN
	// 3    baz
	//
	// datatype: string
}

func ExampleNew_float_precision() {
	s := MustNew([]float64{1.5511, 2.6611})
	fmt.Println(s)
	// Output:
	// 0    1.55
	// 1    2.66
	//
	// datatype: float64
}

func ExampleNew_string_named() {
	s := MustNew([]string{"foo", "bar", "", "baz"}, Config{Name: "foobar"})
	fmt.Println(s)
	// Output:
	// 0    foo
	// 1    bar
	// 2    NaN
	// 3    baz
	//
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
	//
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
	//
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
	//
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
	//
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
	//  0 10000     foo
	//    10100     bar
	//  1 10200     NaN
	//    10300     baz
	//  0 10400     qux
	//    10500    quux
	//
	// datatype: string
	// name: foobar
}

func ExampleNew_repeating_singleIndex() {
	s := MustNew(
		[]string{"foo", "bar", "baz", "qux"},
		Config{Index: []int{0, 0, 1, 1}})
	fmt.Println(s)
	// Output:
	// 0    foo
	//      bar
	// 1    baz
	//      qux
	//
	// datatype: string
}

func ExampleNew_repeating_multiIndex() {
	s := MustNew(
		[]string{"foo", "bar", "baz", "qux"},
		Config{MultiIndex: []interface{}{[]int{0, 0, 1, 1}, []string{"A", "A", "B", "B"}}})
	fmt.Println(s)
	// Output:
	// 0 A    foo
	//        bar
	// 1 B    baz
	//        qux
	//
	// datatype: string
}

func ExampleNew_repeating_allowed() {
	options.SetDisplayRepeatedLabels(true)
	s := MustNew(
		[]string{"foo", "bar", "baz", "qux"},
		Config{Index: []int{0, 0, 1, 1}})
	fmt.Println(s)
	options.RestoreDefaults()
	// Output:
	// 0    foo
	// 0    bar
	// 1    baz
	// 1    qux
	//
	// datatype: string
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
	//
	// datatype: string
	// name: foobar
}

func ExampleNew_datetime_single() {
	s := MustNew([]time.Time{time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)})
	fmt.Println(s)

	// Output:
	// 0    1/1/2019T00:00:00
	//
	// datatype: dateTime
}

func ExampleNew_datetime_manyRows() {
	s := MustNew([]time.Time{
		time.Date(2019, 5, 1, 15, 9, 30, 30, time.UTC),
		time.Date(2019, 5, 2, 15, 15, 55, 55, time.UTC),
	})
	fmt.Println(s)

	// Output:
	// 0    5/1/2019T15:09:30
	// 1    5/2/2019T15:15:55
	//
	// datatype: dateTime
}

func ExampleNew_config_nameOnly() {
	s := MustNew([]string{"foo", "bar"}, Config{Name: "baz"})
	fmt.Println(s)
	// Output:
	// 0    foo
	// 1    bar
	//
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
	//
	// datatype: string
}

func ExampleNew_config_datatype() {
	s := MustNew([]interface{}{"1", "foo"}, Config{DataType: options.Float64})
	fmt.Println(s)
	// Output:
	// 0    1.00
	// 1     NaN
	//
	// datatype: float64
}

func ExampleNew_maxWidth_index() {
	s := MustNew([]string{"foo", "bar"}, Config{Index: []string{"This is a very long index row. Very long indeed.", "qux"}, IndexName: "baz"})
	fmt.Println(s)
	// Output:
	//                                 baz
	// This is a very long index row. V...    foo
	//                                 qux    bar
	//
	// datatype: string
}

func ExampleNew_maxWidth_value() {
	s := MustNew([]string{"This is a very long index row. Very long indeed.", "foo"})
	fmt.Println(s)
	// Output:
	// 0    This is a very long index row. V...
	// 1                                    foo
	//
	//datatype: string
}

func ExampleNew_exceed_maxRows_even() {
	options.SetDisplayMaxRows(2)
	s := MustNew([]float64{0, 1, 2, 3, 4})
	fmt.Println(s)
	options.RestoreDefaults()
	// Output:
	// 0    0.00
	// ...
	// 4    4.00
	//
	//datatype: float64
}

func ExampleNew_exceed_maxRows_odd() {
	options.SetDisplayMaxRows(3)
	s := MustNew([]float64{0, 1, 2, 3, 4})
	fmt.Println(s)
	options.RestoreDefaults()

	// Output:
	// 0    0.00
	// 1    1.00
	// ...
	// 4    4.00
	//
	//datatype: float64
}

// [END constructor examples]

func ExampleSeries_Describe_scalarString() {
	s, _ := New("foo")
	s.name = "bar"
	s.Describe()
	// Output:
	//    len    1
	//  valid    1
	//   null    0
	// unique    1
	//
	// datatype: string
	// name: bar
}

func ExampleSeries_Describe_float() {
	s, _ := New([]float64{1, math.NaN(), 2, 3, 4, 5, math.NaN(), 6, 7, 8, 9})
	s.Describe()
	// Output:
	//   len      11
	// valid       9
	//  null       2
	//  mean    5.00
	//   min    1.00
	//   25%    2.50
	//   50%    5.00
	//   75%    7.50
	//   max    9.00
	//
	// datatype: float64
}

func ExampleSeries_Describe_float_empty() {
	s, _ := New([]float64{})
	s.Describe()
	// Output:
	//   len      0
	// valid      0
	//  null      0
	//  mean    NaN
	//   min    NaN
	//   25%    NaN
	//   50%    NaN
	//   75%    NaN
	//   max    NaN
	//
	// datatype: float64
}

func ExampleSeries_Describe_int() {
	s, _ := New([]int64{1, 2, 3, 4, 5, 6, 7, 8, 9})
	s.Describe()
	// Output:
	//   len       9
	// valid       9
	//  null       0
	//  mean    5.00
	//   min    1.00
	//   25%    2.50
	//   50%    5.00
	//   75%    7.50
	//   max    9.00
	//
	// datatype: int64
}

func ExampleSeries_Describe_string() {
	s, _ := New([]string{"low", "medium", "medium", ""})
	s.Describe()
	// Output:
	//    len    4
	//  valid    3
	//   null    1
	// unique    2
	//
	// datatype: string
}

func ExampleSeries_Describe_string_empty() {
	s, _ := New([]string{"", ""})
	s.Describe()
	// Output:
	//    len    2
	//  valid    0
	//   null    2
	// unique    0
	//
	// datatype: string
}

func ExampleSeries_Describe_bool() {
	s, _ := New([]bool{true, false, false})
	s.Describe()
	// Output:
	//   len       3
	// valid       3
	//  null       0
	//   sum    1.00
	//  mean    0.33
	//
	// datatype: bool
}

func ExampleSeries_Describe_datetime() {
	s, _ := New([]time.Time{
		time.Date(2019, 4, 18, 15, 0, 0, 0, time.UTC),
		time.Date(2019, 4, 19, 15, 0, 0, 0, time.UTC),
		time.Time{},
	})
	s.Describe()
	// Output:
	//      len                                3
	//    valid                                2
	//     null                                1
	//   unique                                2
	// earliest    2019-04-18 15:00:00 +0000 UTC
	//   latest    2019-04-19 15:00:00 +0000 UTC
	//
	// datatype: dateTime
}

func ExampleSeries_Describe_datetime_empty() {
	s, _ := New([]time.Time{time.Time{}})
	s.Describe()
	// Output:
	//      len                                1
	//    valid                                0
	//     null                                1
	//   unique                                0
	// earliest    0001-01-01 00:00:00 +0000 UTC
	//   latest    0001-01-01 00:00:00 +0000 UTC
	//
	// datatype: dateTime
}

func ExampleSeries_Describe_interface() {
	s, _ := New([]interface{}{1.5, 1, "", false})
	s.Describe()
	// Output:
	//   len    4
	// valid    3
	//  null    1
	//
	// datatype: interface
}

// [START additional structs]

func ExampleSeries_empty_series() {
	s := newEmptySeries()
	fmt.Println(s)
	// Output:
	// {Empty Series}
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

func ExampleInPlace_method_list() {
	s := MustNew("foo")
	fmt.Println(s.InPlace)
	// Output:
	// {InPlace Handler}
	// Methods:
	// Append
	// Apply
	// Drop
	// DropDuplicates
	// DropNull
	// DropRows
	// Insert
	// Join
	// Len
	// Less
	// Set
	// SetRows
	// Sort
	// String
	// Subset
	// Swap
	// ToBool
	// ToDateTime
	// ToFloat64
	// ToInt64
	// ToInterface
	// ToString
}

func ExampleIndex_valid_printer() {
	s := MustNew([]string{"foo", "bar", "baz"})
	fmt.Println(s.Index)
	// Output:
	// {Index | Len: 3, NumLevels: 1}
}
func ExampleGrouping_method_list() {
	s := MustNew(
		[]string{"foo", "bar", "baz"},
		Config{MultiIndex: []interface{}{[]int{0, 0, 1}, []int{100, 100, 101}}})
	g := s.GroupByIndex()
	fmt.Println(g)
	// Output:
	// {Grouping | NumGroups: 2, Groups: [0 100, 1 101]}
}
