package series

import (
	"fmt"
	"math"
	"reflect"
)

func (s Series) Sum() (float64, error) {
	switch s.Kind {
	case Float:
		vals := s.Values.(floatValues)
		return vals.sum(), nil
	case Int:
		vals := s.Values.(intValues)
		return vals.sum(), nil
	case Bool:
		vals := s.Values.(boolValues)
		return vals.sum(), nil
	default:
		return math.NaN(), fmt.Errorf("Sum not supported for type %v", s.Kind)
	}
}

func (s Series) Mean() (float64, error) {
	switch s.Kind {
	case Float:
		vals := s.Values.(floatValues)
		return vals.mean(), nil
	case Int:
		vals := s.Values.(intValues)
		return vals.mean(), nil
	default:
		return math.NaN(), fmt.Errorf("Mean not supported for type %v", s.Kind)
	}
}

func (s Series) Median() (float64, error) {
	switch s.Kind {
	case Float:
		vals := s.Values.(floatValues)
		return vals.median(), nil
	case Int:
		vals := s.Values.(intValues)
		return vals.median(), nil
	default:
		return math.NaN(), fmt.Errorf("Median not supported for type %v", s.Kind)
	}
}

func (s Series) Count() int {
	return count(s.Values)
}

func (s Series) String() string {
	switch s.Kind {
	case DateTime:
		var printer string
		vals := s.Values.(dateTimeValues)
		for _, val := range vals {
			printer += fmt.Sprintln(val.v.Format("01/02/2006"))
		}
		return printer
	default:
		return print(s.Values)
	}
}

// expects to receive a slice of typed value structs (eg pd.floatValues)
// each struct must contain a boolean field called "null"
func count(data interface{}) int {
	var count int
	d := reflect.ValueOf(data)
	for i := 0; i < d.Len(); i++ {
		if !d.Index(i).FieldByName("null").Bool() {
			count++
		}
	}
	return count
}

// expects to receive a slice of typed value structs (eg pd.floatValues)
// each struct must contain a boolean field called "v"
func print(data interface{}) string {
	d := reflect.ValueOf(data)
	var printer string
	for i := 0; i < d.Len(); i++ {
		printer += fmt.Sprintln(d.Index(i).FieldByName("v"))
	}
	return printer
}
