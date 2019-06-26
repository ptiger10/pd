package series

import (
	"math"
	"sort"

	gonum "github.com/gonum/floats" // uses optimized gonum/floats methods where available
	"github.com/montanaflynn/stats" // and stats package otherwise
	"github.com/ptiger10/pd/options"
)

// an interface of valid (non-null) values; appropriate for type assertion
func (s *Series) validVals() interface{} {
	valid, _ := s.values.Subset(s.valid())
	return valid.Vals()
}

// Sum of non-null series elements. For bool values, sum of true values.
//
// Applies to: Float, Int, Bool. If inapplicable, defaults to math.Nan().
func (s *Series) Sum() float64 {
	vals := s.validVals()
	switch s.datatype {
	case options.Float64, options.Int64:
		data := ensureFloatFromNumerics(vals)
		return gonum.Sum(data)
	case options.Bool:
		var sum float64
		data := ensureBools(vals)
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
//
// Applies to: Float, Int. If inapplicable, defaults to math.Nan().
func (s *Series) Mean() float64 {
	var sum float64
	vals := s.validVals()
	switch s.datatype {
	case options.Float64, options.Int64:
		data := ensureFloatFromNumerics(vals)
		for _, d := range data {
			sum += d
		}
		return sum / float64(len(data))
	case options.Bool:
		data := ensureBools(vals)
		l := len(data)
		return s.Sum() / float64(l)
	default:
		return math.NaN()
	}
}

// Median of a series.
//
// Applies to: Float, Int. If inapplicable, defaults to math.Nan().
func (s *Series) Median() float64 {
	vals := s.validVals()
	switch s.datatype {
	case options.Float64, options.Int64:
		data := ensureFloatFromNumerics(vals)
		if len(data) == 0 {
			return math.NaN()
		}
		sort.Float64s(data)
		mNumber := len(data) / 2
		if len(data)%2 != 0 { // checks if sequence has odd number of elements
			return data[mNumber]
		}
		return (data[mNumber-1] + data[mNumber]) / 2
	default:
		return math.NaN()
	}
}

// Min of a series.
//
// Applies to: Float, Int. If inapplicable, defaults to math.Nan().
func (s *Series) Min() float64 {
	vals := s.validVals()
	switch s.datatype {
	case options.Float64, options.Int64:
		data := ensureFloatFromNumerics(vals)
		if len(data) == 0 {
			return math.NaN()
		}
		return gonum.Min(data)
	default:
		return math.NaN()
	}
}

// Max of a series.
//
// Applies to: Float, Int. If inapplicable, defaults to math.Nan().
func (s *Series) Max() float64 {
	vals := s.validVals()
	switch s.datatype {
	case options.Float64, options.Int64:
		data := ensureFloatFromNumerics(vals)
		if len(data) == 0 {
			return math.NaN()
		}
		return gonum.Max(data)
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
	vals := s.validVals()
	switch s.datatype {
	case options.Float64, options.Int64:
		data := ensureFloatFromNumerics(vals)
		if len(data) == 0 {
			return math.NaN()
		}
		std, err := stats.StandardDeviation(data)
		if err != nil {
			return math.NaN()
		}
		return std
	default:
		return math.NaN()
	}
}
