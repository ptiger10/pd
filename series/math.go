package series

import (
	"math"
	"sort"

	gonum "github.com/gonum/floats" // uses optimized gonum/floats methods where available
	"github.com/montanaflynn/stats" // and stats package otherwise
	"github.com/ptiger10/pd/kinds"
)

// Math contains mathematical methods
type Math struct {
	s *Series
}

// an interface of valid (non-null) values; appropriate for type assertion
func (s Series) validVals() interface{} {
	valid, _ := s.values.In(s.valid())
	return valid.Vals()
}

// Sum is shorthand for series.Math.Sum()
//
// Applies to: Float, Int, Bool. If inapplicable, defaults to math.Nan().
func (s Series) Sum() float64 {
	return s.Math.Sum()
}

// Mean is shorthand for series.Math.Mean()
//
// Applies to: Float, Int, Bool. If inapplicable, defaults to math.Nan().
func (s Series) Mean() float64 {
	return s.Math.Mean()
}

// Sum of non-null series elements. For bool values, sum of true values.
//
// Applies to: Float, Int, Bool. If inapplicable, defaults to math.Nan().
func (m Math) Sum() float64 {
	vals := m.s.validVals()
	switch m.s.kind {
	case kinds.Float64, kinds.Int64:
		data := ensureFloatFromNumerics(vals)
		return gonum.Sum(data)
	case kinds.Bool:
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
func (m Math) Mean() float64 {
	var sum float64
	vals := m.s.validVals()
	switch m.s.kind {
	case kinds.Float64, kinds.Int64:
		data := ensureFloatFromNumerics(vals)
		for _, d := range data {
			sum += d
		}
		return sum / float64(len(data))
	case kinds.Bool:
		data := ensureBools(vals)
		l := len(data)
		return m.Sum() / float64(l)
	default:
		return math.NaN()
	}
}

// Median of a series.
//
// Applies to: Float, Int. If inapplicable, defaults to math.Nan().
func (m Math) Median() float64 {
	vals := m.s.validVals()
	switch m.s.kind {
	case kinds.Float64, kinds.Int64:
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
func (m Math) Min() float64 {
	vals := m.s.validVals()
	switch m.s.kind {
	case kinds.Float64, kinds.Int64:
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
func (m Math) Max() float64 {
	vals := m.s.validVals()
	switch m.s.kind {
	case kinds.Float64, kinds.Int64:
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
func (m Math) Quartile(i int) float64 {
	vals := m.s.validVals()
	switch m.s.kind {
	case kinds.Float64, kinds.Int64:
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
func (m Math) Std() float64 {
	vals := m.s.validVals()
	switch m.s.kind {
	case kinds.Float64, kinds.Int64:
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
