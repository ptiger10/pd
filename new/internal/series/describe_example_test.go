package series

import (
	"math"

	"github.com/ptiger10/pd/series"
)

func ExampleSeries_Describe_float() {
	s, _ := New([]float64{1, math.NaN(), 2, 3, 4, 5, math.NaN(), 6, 7, 8, 9})
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
