package series_test

import (
	"fmt"
	"math"
	"time"

	"github.com/ptiger10/pd/series"
)

func ExampleSeries_default() {
	s, _ := series.New([]int64{1, 3, 5})
	fmt.Println(s)
	// Output:
	// 1
	// 3
	// 5
}
func ExampleSeries_dateTime() {
	s, _ := series.New([]time.Time{
		time.Date(2019, 4, 18, 15, 0, 0, 0, time.UTC),
		time.Date(2019, 4, 19, 15, 0, 0, 0, time.UTC)})

	fmt.Println(s)
	// Output:
	// 04/18/2019
	// 04/19/2019

}

func ExampleSeries_Describe_float() {
	s, _ := series.New([]float64{1.0, math.NaN(), 3.0, 5.0})
	s.Describe()
	// Output:
	// len     4
	// valid   3
	// null    1
	// mean    3.0000
	// median  3.0000
	// min     1.0000
	// max     5.0000
}

func ExampleSeries_Describe_int() {
	s, _ := series.New([]int64{-1, 0, 1})
	s.Describe()
	// Output:
	// len     3
	// valid   3
	// null    0
	// mean    0.0000
	// median  0.0000
	// min     -1.0000
	// max     1.0000
}

func ExampleSeries_Describe_string() {
	s, _ := series.New([]string{"low", "medium", "medium", ""})
	s.Describe()
	// Output:
	// len     4
	// valid   3
	// null    1
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
