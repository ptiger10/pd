package series

import (
	"strings"
	"time"
)

func (s *Series) Filter(cmp func(interface{}) bool) []int {
	vals := s.all()
	include := make([]int, 0)
	for i, val := range vals {
		if cmp(val) {
			include = append(include, i)
		}
	}
	return include
}

// customFilterFloat64 converts a Series to float values, applies a filter, and returns the rows where the condition is true.
func (s *Series) customFilterFloat64(cmp func(float64) bool) []int {
	include := make([]int, 0)
	vals := s.ToFloat64().values.Vals().([]float64)
	for i, val := range vals {
		if cmp(val) {
			include = append(include, i)
		}
	}
	return include
}

// customFilterString converts a Series to string values, applies a filter, and returns the rows where the condition is true.
func (s *Series) customFilterString(cmp func(string) bool) []int {
	include := make([]int, 0)
	vals := s.ToString().values.Vals().([]string)
	for i, val := range vals {
		if cmp(val) {
			include = append(include, i)
		}
	}
	return include
}

// customFilterBool converts a Series to bool values, applies a filter, and returns the rows where the condition is true.
func (s *Series) customFilterBool(cmp func(bool) bool) []int {
	include := make([]int, 0)
	vals := s.ToBool().values.Vals().([]bool)
	for i, val := range vals {
		if cmp(val) {
			include = append(include, i)
		}
	}
	return include
}

// customFilterDateTime converts a Series to datetime values, applies a filter, and returns the rows where the condition is true.
func (s *Series) customFilterDateTime(cmp func(time.Time) bool) []int {
	include := make([]int, 0)
	vals := s.ToDateTime().values.Vals().([]time.Time)
	for i, val := range vals {
		if cmp(val) {
			include = append(include, i)
		}
	}
	return include
}

// GT filter: Greater Than (numeric).
func (s *Series) GT(comparison float64) []int {
	return s.customFilterFloat64(func(elem float64) bool {
		return elem > comparison
	})
}

// GTE filter: Greater Than or Equal To (numeric).
func (s *Series) GTE(comparison float64) []int {
	return s.customFilterFloat64(func(elem float64) bool {
		return elem >= comparison
	})
}

// LT filter - Less Than (numeric).
func (s *Series) LT(comparison float64) []int {
	return s.customFilterFloat64(func(elem float64) bool {
		return elem < comparison
	})
}

// LTE filter - Less Than or Equal To (numeric).
func (s *Series) LTE(comparison float64) []int {
	return s.customFilterFloat64(func(elem float64) bool {
		return elem <= comparison
	})
}

// EQ filter - Equal To (numeric).
func (s *Series) EQ(comparison float64) []int {
	return s.customFilterFloat64(func(elem float64) bool {
		return elem == comparison
	})
}

// NEQ filter - Not Equal To (numeric).
func (s *Series) NEQ(comparison float64) []int {
	return s.customFilterFloat64(func(elem float64) bool {
		return elem != comparison
	})
}

// Contains filter - value contains substr (string).
func (s *Series) Contains(substr string) []int {
	return s.customFilterString(func(elem string) bool {
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
	return s.customFilterBool(func(elem bool) bool {
		return elem
	})
}

// False filter - value is false (bool).
func (s *Series) False() []int {
	return s.customFilterBool(func(elem bool) bool {
		return !elem
	})
}

// Before filter - value is before a specific time (datetime).
func (s *Series) Before(t time.Time) []int {
	return s.customFilterDateTime(func(elem time.Time) bool {
		return elem.Before(t)
	})
}

// After filter - value is after a specific time (datetime).
func (s *Series) After(t time.Time) []int {
	return s.customFilterDateTime(func(elem time.Time) bool {
		return elem.After(t)
	})
}
