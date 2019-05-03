package series

import (
	"log"
	"math"
	"sort"

	gonum "github.com/gonum/floats" // uses optimized gonum/floats methods where available
	"github.com/montanaflynn/stats" // and stats package otherwise
	"github.com/ptiger10/pd/new/kinds"
)

// appropriate for numeric data only ([]float64 or []int64)
func ensureFloatFromNumerics(vals interface{}) []float64 {
	var data []float64
	if ints, ok := vals.([]int64); ok {
		data = convertIntToFloat(ints)
	} else if floats, ok := vals.([]float64); ok {
		data = floats
	} else {
		log.Printf("EnsureFloatFromNumerics has received an unallowable value: %v", vals)
		return nil
	}
	return data
}

func convertIntToFloat(vals []int64) []float64 {
	var ret []float64
	for _, val := range vals {
		ret = append(ret, float64(val))
	}
	return ret
}

// Sum of a series.
//
// Applies to: Float, Int
func (s Series) Sum() float64 {
	vals := s.validVals()
	switch s.Kind {
	case kinds.Float, kinds.Int:
		data := ensureFloatFromNumerics(vals)
		return gonum.Sum(data)
	default:
		return math.NaN()
	}
}

// Mean of a series.
//
// Applies to: Float, Int
func (s Series) Mean() float64 {
	var sum float64
	vals := s.validVals()
	switch s.Kind {
	case kinds.Float, kinds.Int:
		data := ensureFloatFromNumerics(vals)
		for _, d := range data {
			sum += d
		}
		return sum / float64(len(data))
	default:
		return math.NaN()
	}
}

// Median of a series.
//
// Applies to: Float, Int
func (s Series) Median() float64 {
	vals := s.validVals()
	switch s.Kind {
	case kinds.Float, kinds.Int:
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
// Applies to: Float, Int
func (s Series) Min() float64 {
	vals := s.validVals()
	switch s.Kind {
	case kinds.Float, kinds.Int:
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
// Applies to: Float, Int
func (s Series) Max() float64 {
	vals := s.validVals()
	switch s.Kind {
	case kinds.Float, kinds.Int:
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
// Applies to: Float, Int
func (s Series) Quartile(i int) float64 {
	vals := s.validVals()
	switch s.Kind {
	case kinds.Float, kinds.Int:
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
// Applies to: Float, Int
func (s Series) Std() float64 {
	vals := s.validVals()
	switch s.Kind {
	case kinds.Float, kinds.Int:
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
