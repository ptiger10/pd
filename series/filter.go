package series

import (
	"strings"
	"time"
)

// Apply a callback function to every value in a Series and return a new Series.
// The Apply function iterates over all Series values in interface{} form and applies the callback function to each.
// The final values are then converted to match the datatype of the original Series.
// The caller is responsible for handling the type assertion on the interface, though this step is not necessary if the datatype is known with certainty.
// For example, here are two ways to write an apply function that computes the z-score of every row and rounds it two decimal points:
//
// #1 (safer) error check type assertion
//
//  s.Apply(func(val interface{}) interface{} {
// 		v, ok := val.(float64)
// 			if !ok {
//				return ""
// 			}
// 		return (v - s.Mean()) / s.Std()
//
// Input:
// 0    1
// 1    2
// 2    3
//
// Output:
// 0    -1.22...
// 1        0
// 2     1.22...
//
// #2 (riskier) no error check
//
//  s.Apply(func(val interface{}) interface{} {
// 		return (val.(float64) - s.Mean()) / s.Std()
// 	})
func (s *Series) Apply(fn func(interface{}) interface{}) *Series {
	vals := s.all()
	newVals := make([]interface{}, 0)
	for _, val := range vals {
		newVal := fn(val)
		newVals = append(newVals, newVal)
	}
	// ducks error because []interface{} as arg in New constructor cannot trigger unsupported error
	ret, _ := New(newVals, Config{DataType: s.datatype})
	ret.index = s.index

	return ret
}

// Filter a Series using a callback function test.
// The Filter function iterates over all Series values in interface{} form and applies the callback test to each.
// The return value is a slice of integer positions of all the rows passing the test.
// The caller is responsible for handling the type assertion on the interface, though this step is not necessary if the datatype is known with certainty.
// For example, here are two ways to write a filter that returns all rows with the suffix "boo":
//
// #1 (safer) error check type assertion
//
//  s.Filter(func(val interface{}) bool {
//		v, ok := val.(string)
//		if !ok {
// 			return false
//		}
//		if strings.HasSuffix(v, "boo") {
// 			return true
// 		}
// 		return false
// 	})
//
// Input:
// 0    bamboo
// 1    leaves
// 2    taboo
//
// Output:
// []int{0,2}
//
// #2 (riskier) no error check
//
//  s.Filter(func(val interface{}) bool {
//		if strings.HasSuffix(val.(string), "boo") {
// 			return true
// 		}
// 		return false
// 	})
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

// filterFloat64 converts a Series to float values, applies a filter, and returns the rows where the condition is true.
func (s *Series) filterFloat64(cmp func(float64) bool) []int {
	include := make([]int, 0)
	vals := s.ToFloat64().values.Vals().([]float64)
	for i, val := range vals {
		if cmp(val) {
			include = append(include, i)
		}
	}
	return include
}

// filterString converts a Series to string values, applies a filter, and returns the rows where the condition is true.
func (s *Series) filterString(cmp func(string) bool) []int {
	include := make([]int, 0)
	vals := s.ToString().values.Vals().([]string)
	for i, val := range vals {
		if cmp(val) {
			include = append(include, i)
		}
	}
	return include
}

// filterBool converts a Series to bool values, applies a filter, and returns the rows where the condition is true.
func (s *Series) filterBool(cmp func(bool) bool) []int {
	include := make([]int, 0)
	vals := s.ToBool().values.Vals().([]bool)
	for i, val := range vals {
		if cmp(val) {
			include = append(include, i)
		}
	}
	return include
}

// filterDateTime converts a Series to datetime values, applies a filter, and returns the rows where the condition is true.
func (s *Series) filterDateTime(cmp func(time.Time) bool) []int {
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
	return s.filterFloat64(func(elem float64) bool {
		return elem > comparison
	})
}

// GTE filter: Greater Than or Equal To (numeric).
func (s *Series) GTE(comparison float64) []int {
	return s.filterFloat64(func(elem float64) bool {
		return elem >= comparison
	})
}

// LT filter - Less Than (numeric).
func (s *Series) LT(comparison float64) []int {
	return s.filterFloat64(func(elem float64) bool {
		return elem < comparison
	})
}

// LTE filter - Less Than or Equal To (numeric).
func (s *Series) LTE(comparison float64) []int {
	return s.filterFloat64(func(elem float64) bool {
		return elem <= comparison
	})
}

// EQ filter - Equal To (numeric).
func (s *Series) EQ(comparison float64) []int {
	return s.filterFloat64(func(elem float64) bool {
		return elem == comparison
	})
}

// NEQ filter - Not Equal To (numeric).
func (s *Series) NEQ(comparison float64) []int {
	return s.filterFloat64(func(elem float64) bool {
		return elem != comparison
	})
}

// Contains filter - value contains substr (string).
func (s *Series) Contains(substr string) []int {
	return s.filterString(func(elem string) bool {
		return strings.Contains(elem, substr)
	})
}

// InList filter - value is contained within list (string).
func (s *Series) InList(list []string) []int {
	return s.filterString(func(elem string) bool {
		for _, s := range list {
			if elem == s {
				return true
			}
		}
		return false
	})
}

// True filter - value is true (bool).
func (s *Series) True() []int {
	return s.filterBool(func(elem bool) bool {
		return elem
	})
}

// False filter - value is false (bool).
func (s *Series) False() []int {
	return s.filterBool(func(elem bool) bool {
		return !elem
	})
}

// Before filter - value is before a specific time (datetime).
func (s *Series) Before(t time.Time) []int {
	return s.filterDateTime(func(elem time.Time) bool {
		return elem.Before(t)
	})
}

// After filter - value is after a specific time (datetime).
func (s *Series) After(t time.Time) []int {
	return s.filterDateTime(func(elem time.Time) bool {
		return elem.After(t)
	})
}
