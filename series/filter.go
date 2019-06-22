package series

import (
	"log"
	"strings"
	"time"

	"github.com/ptiger10/pd/options"
)

// Subset returns a subset of a Series based on the supplied integer positions.
func (s *Series) Subset(rows []int) *Series {
	sub, err := s.in(rows)
	if err != nil {
		if options.GetLogWarnings() {
			log.Printf("s.Subset(): %v", err)
		}
		return newEmptySeries()
	}
	return sub
}

// CustomFilterFloat64 converts a Series to float values, applies a filter, and returns the rows where the condition is true.
func (s *Series) CustomFilterFloat64(cmp func(float64) bool) []int {
	include := make([]int, 0)
	vals := s.ToFloat64().values.Vals().([]float64)
	for i, val := range vals {
		if cmp(val) {
			include = append(include, i)
		}
	}
	return include
}

// CustomFilterString converts a Series to string values, applies a filter, and returns the rows where the condition is true.
func (s *Series) CustomFilterString(cmp func(string) bool) []int {
	include := make([]int, 0)
	vals := s.ToString().values.Vals().([]string)
	for i, val := range vals {
		if cmp(val) {
			include = append(include, i)
		}
	}
	return include
}

// CustomFilterBool converts a Series to bool values, applies a filter, and returns the rows where the condition is true.
func (s *Series) CustomFilterBool(cmp func(bool) bool) []int {
	include := make([]int, 0)
	vals := s.ToBool().values.Vals().([]bool)
	for i, val := range vals {
		if cmp(val) {
			include = append(include, i)
		}
	}
	return include
}

// CustomFilterDateTime converts a Series to datetime values, applies a filter, and returns the rows where the condition is true.
func (s *Series) CustomFilterDateTime(cmp func(time.Time) bool) []int {
	include := make([]int, 0)
	vals := s.ToDateTime().values.Vals().([]time.Time)
	for i, val := range vals {
		if cmp(val) {
			include = append(include, i)
		}
	}
	return include
}

// Gt filter: Greater Than (numeric).
func (s *Series) Gt(comparison float64) []int {
	return s.CustomFilterFloat64(func(elem float64) bool {
		return elem > comparison
	})
}

// Gte filter: Greater Than or Equal To (numeric).
func (s *Series) Gte(comparison float64) []int {
	return s.CustomFilterFloat64(func(elem float64) bool {
		return elem >= comparison
	})
}

// Lt filter - Less Than (numeric).
func (s *Series) Lt(comparison float64) []int {
	return s.CustomFilterFloat64(func(elem float64) bool {
		return elem < comparison
	})
}

// Lte filter - Less Than or Equal To (numeric).
func (s *Series) Lte(comparison float64) []int {
	return s.CustomFilterFloat64(func(elem float64) bool {
		return elem <= comparison
	})
}

// Eq filter - Equal To (numeric).
func (s *Series) Eq(comparison float64) []int {
	return s.CustomFilterFloat64(func(elem float64) bool {
		return elem == comparison
	})
}

// Neq filter - Not Equal To (numeric).
func (s *Series) Neq(comparison float64) []int {
	return s.CustomFilterFloat64(func(elem float64) bool {
		return elem != comparison
	})
}

// Contains filter - value contains substr (string).
func (s *Series) Contains(substr string) []int {
	return s.CustomFilterString(func(elem string) bool {
		return strings.Contains(elem, substr)
	})
}

// // In filter - value is contained within list (string).
// func (s *Series) In(list []string) []int {
// 	return func(elem string) bool {
// 		for _, s := range list {
// 			if elem == s {
// 				return true
// 			}
// 		}
// 		return false
// 	}
// }

// True filter - value is true (bool).
func (s *Series) True() []int {
	return s.CustomFilterBool(func(elem bool) bool {
		return elem
	})
}

// False filter - value is false (bool).
func (s *Series) False() []int {
	return s.CustomFilterBool(func(elem bool) bool {
		return !elem
	})
}

// Before filter - value is before a specific time (datetime).
func (s *Series) Before(t time.Time) []int {
	return s.CustomFilterDateTime(func(elem time.Time) bool {
		return elem.Before(t)
	})
}

// After filter - value is after a specific time (datetime).
func (s *Series) After(t time.Time) []int {
	return s.CustomFilterDateTime(func(elem time.Time) bool {
		return elem.After(t)
	})
}
