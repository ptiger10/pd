package series_test

import (
	"fmt"
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
