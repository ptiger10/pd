package series

import (
	"math"
	"runtime"
	"sort"
	"sync"

	// uses optimized gonum/floats methods where available
	"github.com/montanaflynn/stats" // and stats package otherwise
	"github.com/ptiger10/pd/options"
)

// an interface of valid (non-null) values; appropriate for type assertion
func (s *Series) validVals() interface{} {
	valid := s.values.Subset(s.valid())
	return valid.Vals()
}

// Sum of non-null float64 or int64 Series values. For bool values, sum of true values. If inapplicable, defaults to math.Nan().
func (s *Series) Sum() float64 {
	var sum float64
	var wg sync.WaitGroup

	// null int values are represented as 0, but that's ok for sum
	calc := func(data []float64) float64 {
		var sum float64
		for _, d := range data {
			if !math.IsNaN(d) {
				sum += d
			}
		}
		return sum
	}

	awaitSum := func(data []float64, ch chan<- float64) {
		sum := calc(data)
		ch <- sum
		wg.Done()
	}

	switch s.datatype {
	case options.Float64, options.Int64:
		data := ensureFloatFromNumerics(s.Vals())

		if !options.GetAsync() {
			return calc(data)
		}
		numPartitions := runtime.GOMAXPROCS(0)
		valsPerPartition := len(data) / numPartitions
		ch := make(chan float64, numPartitions)
		for i := 0; i < numPartitions; i++ {
			var sub []float64
			if i != numPartitions-1 {
				sub = data[i*valsPerPartition : (i+1)*(valsPerPartition)]
			} else {
				sub = data[i*valsPerPartition:]
			}
			wg.Add(1)
			go awaitSum(sub, ch)
		}
		wg.Wait()
		close(ch)
		for partitionSum := range ch {
			sum += partitionSum
		}
		return sum

	case options.Bool:
		data := ensureBools(s.Vals())
		// null bool values are represented as false, but that's ok for sum
		for _, d := range data {
			if d {
				sum++
			}
		}
		return sum
	default:
		return math.NaN()
	}
}

// Mean of non-null series values. For bool values, mean of true values.
// Applies to float64 and int64. If inapplicable, defaults to math.Nan().
func (s *Series) Mean() float64 {
	switch s.datatype {
	case options.Float64:
		var sum float64
		var counter int
		data := ensureFloatFromNumerics(s.Vals())
		for _, d := range data {
			if !math.IsNaN(d) {
				sum += d
				counter++
			}
		}
		return sum / float64(counter)
	case options.Int64:
		return s.Sum() / float64(s.validCount())
	case options.Bool:
		return s.Sum() / float64(s.validCount())
	default:
		return math.NaN()
	}
}

// Median of a series. Applies to float64 and int64. If inapplicable, defaults to math.Nan().
func (s *Series) Median() float64 {
	calculateMedian := func(data []float64) float64 {
		if len(data) == 0 {
			return math.NaN()
		}
		// rounds down if there are even number of elements
		mNumber := len(data) / 2

		// odd number of elements
		if len(data)%2 != 0 {
			return data[mNumber]
		}
		// even number of elements
		return (data[mNumber-1] + data[mNumber]) / 2
	}

	switch s.datatype {
	case options.Float64:
		// sort then slice out the nulls from the beginning
		data := ensureFloatFromNumerics(s.Vals())
		sort.Float64s(data)
		numValids := s.validCount()
		firstValid := len(data) - numValids
		data = data[firstValid:]
		return calculateMedian(data)

	case options.Int64:
		vals := s.validVals()
		data := ensureFloatFromNumerics(vals)
		sort.Float64s(data)
		return calculateMedian(data)

	default:
		return math.NaN()
	}
}

// Min of a series. Applies to float64 and int64. If inapplicable, defaults to math.Nan().
func (s *Series) Min() float64 {
	min := math.NaN()
	switch s.datatype {
	case options.Float64:
		data := ensureFloatFromNumerics(s.Vals())
		for _, d := range data {
			if !math.IsNaN(d) {
				if math.IsNaN(min) {
					min = d
				} else if d < min {
					min = d
				}
			}
		}
		return min
	case options.Int64:
		data := ensureFloatFromNumerics(s.Vals())
		for i, d := range data {
			if !s.values.Null(i) {
				if math.IsNaN(min) {
					min = d
				} else if d < min {
					min = d
				}
			}
		}
		return min
	default:
		return math.NaN()
	}
}

// Max of a series. Applies to float64 and int64. If inapplicable, defaults to math.Nan().
func (s *Series) Max() float64 {
	max := math.NaN()
	switch s.datatype {
	case options.Float64:
		data := ensureFloatFromNumerics(s.Vals())
		for _, d := range data {
			if !math.IsNaN(d) {
				if math.IsNaN(max) {
					max = d
				} else if d > max {
					max = d
				}
			}
		}
		return max
	case options.Int64:
		data := ensureFloatFromNumerics(s.Vals())
		for i, d := range data {
			if !s.values.Null(i) {
				if math.IsNaN(max) {
					max = d
				} else if d > max {
					max = d
				}
			}
		}
		return max
	default:
		return math.NaN()
	}
}

// Quartile i (must be 1, 2, 3)
//
// Applies to: Float, Int. If inapplicable, defaults to math.Nan().
func (s *Series) Quartile(i int) float64 {
	vals := s.validVals()
	switch s.datatype {
	case options.Float64, options.Int64:
		data := ensureFloatFromNumerics(vals)
		val, err := stats.Quartile(data)
		if err != nil {
			return math.NaN()
		}
		switch i {
		case 1:
			return val.Q1
		case 2:
			return val.Q2
		case 3:
			return val.Q3
		default:
			return math.NaN()
		}
	default:
		return math.NaN()
	}
}

// Std returns the Standard Deviation of a series.
//
// Applies to: Float, Int. If inapplicable, defaults to math.Nan().
func (s *Series) Std() float64 {
	switch s.datatype {
	case options.Float64, options.Int64:
		data := ensureFloatFromNumerics(s.Vals())
		mean := s.Mean()
		var variance float64
		var counter int
		for i, d := range data {
			if !s.values.Null(i) {
				variance += (d - mean) * (d - mean)
				counter++
			}
		}
		return math.Pow(variance/float64(counter), 0.5)
	default:
		return math.NaN()
	}
}
