package series_test

import (
	"fmt"
	"math"
	"time"

	"github.com/ptiger10/pd/series"
)

func ExampleSeries_int_printer() {
	s, _ := series.New([]int64{1, 3, 5})
	fmt.Println(s)
	// Output:
	// 1
	// 3
	// 5
}
func ExampleSeries_dateTime_printer() {
	s, _ := series.New([]time.Time{
		time.Date(2019, 4, 18, 15, 0, 0, 0, time.UTC),
		time.Date(2019, 4, 19, 15, 0, 0, 0, time.UTC)})

	fmt.Println(s)
	// Output:
	// 04/18/2019
	// 04/19/2019

}

func ExampleSeries_Describe_float() {
	s, _ := series.New([]float64{1, math.NaN(), 2, 3, 4, 5, math.NaN(), 6, 7, 8, 9})
	s.Describe()
	// Output:
	// len     11
	// valid   9
	// null    2
	// mean    5.0000
	// min     1.0000
	// 25%     2.5000
	// 50%     5.0000
	// 75%     7.5000
	// max     9.0000
}

func ExampleSeries_Describe_float_empty() {
	s, _ := series.New([]float32{})
	s.Describe()
	// Output:
	// len     0
	// valid   0
	// null    0
	// mean    NaN
	// min     NaN
	// 25%     NaN
	// 50%     NaN
	// 75%     NaN
	// max     NaN
}

func ExampleSeries_Describe_int() {
	s, _ := series.New([]int64{1, 2, 3, 4, 5, 6, 7, 8, 9})
	s.Describe()
	// Output:
	// len     9
	// valid   9
	// null    0
	// mean    5.0000
	// min     1.0000
	// 25%     2.5000
	// 50%     5.0000
	// 75%     7.5000
	// max     9.0000
}

func ExampleSeries_Describe_int_empty() {
	s, _ := series.New([]interface{}{"", ""}, series.Type(series.Int))
	s.Describe()
	// Output:
	// len     2
	// valid   0
	// null    2
	// mean    NaN
	// min     NaN
	// 25%     NaN
	// 50%     NaN
	// 75%     NaN
	// max     NaN
}

func ExampleSeries_Describe_string() {
	s, _ := series.New([]string{"low", "medium", "medium", ""})
	s.Describe()
	// Output:
	// len     4
	// valid   3
	// null    1
	// unique  2
}

func ExampleSeries_Describe_string_empty() {
	s, _ := series.New([]string{"", ""})
	s.Describe()
	// Output:
	// len     2
	// valid   0
	// null    2
	// unique  0
}

func ExampleSeries_Describe_bool() {
	s, _ := series.New([]bool{true, false, false})
	s.Describe()
	// Output:
	// len     3
	// valid   3
	// null    0
	// sum     1
}

func ExampleSeries_Describe_datetime() {
	s, _ := series.New([]time.Time{
		time.Date(2019, 4, 18, 15, 0, 0, 0, time.UTC),
		time.Date(2019, 4, 19, 15, 0, 0, 0, time.UTC),
		time.Time{}})
	s.Describe()
	// Output:
	// len     3
	// valid   2
	// null    1
}

func ExampleSeries_Describe_datetime_empty() {
	s, _ := series.New([]time.Time{time.Time{}})
	s.Describe()
	// Output:
	// len     1
	// valid   0
	// null    1
}
