package series

import "fmt"

// Apply contains methods for element-wise manipulation
type Apply struct {
	s *Series
}

// Float64 applies a function to all float64 values and returns a new Series.
func (a Apply) Float64(fn func(float64) float64) (*Series, error) {
	s := a.s.Copy()
	vals, ok := s.values.Vals().([]float64)
	if !ok {
		return nil, fmt.Errorf("float64 Apply expects float64 values only, got %v", a.s.DataType())
	}
	for i := 0; i < a.s.values.Len(); i++ {
		a.s.values.Set(i, fn(vals[i]))
	}
	return s, nil
}

// Multiply multiplies every element by factor.
func (a Apply) Multiply(factor float64) (*Series, error) {
	s, err := a.Float64(func(elem float64) float64 {
		return elem * factor
	})
	if err != nil {
		return nil, fmt.Errorf("Apply.Multiply(): %v", err)
	}
	return s, nil
}
