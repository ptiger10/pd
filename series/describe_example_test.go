package series

import (
	"fmt"
	"math"
	"time"
)

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

func ExampleIndex_valid_printer() {
	s := MustNew([]string{"foo", "bar", "baz"})
	fmt.Println(s.Index)
	// Output:
	// {Index | Len: 3, NumLevels: 1}
}
