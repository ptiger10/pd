package series

import (
	"fmt"
	"time"
)

// Filter contains filtering methods.
type Filter struct {
	s *Series
}

// Float64 filters
func (f Filter) Float64(cmp func(float64) bool) (Series, error) {
	s, _ := f.s.in(f.s.valid())
	var include []int
	vals, ok := s.values.Vals().([]float64)
	if !ok {
		return Series{}, fmt.Errorf("float64 filter expects float64 values only, got %T", f.s.Kind())
	}
	for i, val := range vals {
		if cmp(val) {
			include = append(include, i)
		}
	}
	return s.in(include)
}

// Gt - Greater Than
//
// Applies to: Float, Int
func (f Filter) Gt(comparison float64) (Series, error) {
	s, err := f.Float64(func(elem float64) bool {
		return elem > comparison
	})
	if err != nil {
		return Series{}, fmt.Errorf("Filter.Gt(): %v", err)
	}
	return s, nil
}

// Gt - Greater Than
//
// Applies to: Float, Int
// func (f Filter) Gt(comparison float64) func(float64) bool {
// 	return func(elem float64) bool {
// 		return elem > comparison
// 	}
// }

// Gt - Greater Than
//
// Applies to: Float, Int
func Gt(comparison float64) func(float64) bool {
	return func(elem float64) bool {
		return elem > comparison
	}
}

// Gte - Greater Than or Equal To
//
// Applies to: Float, Int
func Gte(comparison float64) func(float64) bool {
	return func(elem float64) bool {
		return elem >= comparison
	}
}

// Lt - Less Than
//
// Applies to: Float, Int
func Lt(comparison float64) func(float64) bool {
	return func(elem float64) bool {
		return elem < comparison
	}
}

// Lte - Less Than or Equal To
//
// Applies to: Float, Int
func Lte(comparison float64) func(float64) bool {
	return func(elem float64) bool {
		return elem <= comparison
	}
}

// Eq - Equal To
//
// Applies to: Float, Int
func Eq(comparison float64) func(float64) bool {
	return func(elem float64) bool {
		return elem == comparison
	}
}

// In - Contained within slice of strings
//
// Applies to: String
func In(list []string) func(string) bool {
	return func(elem string) bool {
		for _, s := range list {
			if elem == s {
				return true
			}
		}
		return false
	}
}

// True - True, non-null value
//
// Applies to: Bool
func True() func(bool) bool {
	return func(elem bool) bool {
		return elem
	}
}

// Before - Before a specific datetime
//
// Applies to time.Time
func Before(t time.Time) func(time.Time) bool {
	return func(elem time.Time) bool {
		return elem.Before(t)
	}
}

// On - On a specific datetime
//
// Applies to time.Time
func On(t time.Time) func(time.Time) bool {
	return func(elem time.Time) bool {
		return elem.Equal(t)
	}
}

// After - After a specific datetime
//
// Applies to time.Time
func After(t time.Time) func(time.Time) bool {
	return func(elem time.Time) bool {
		return elem.After(t)
	}
}
