package series

import (
	"fmt"
	"log"
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
		return math.NaN(), fmt.Errorf("Sum not supported for Series type %v", s.Kind)
	}
}

func (s Series) Describe() {
	fmt.Println(s.Values.describe())
}

func (s Series) AddConst(c interface{}) (Series, error) {
	switch s.Kind {
	case Float:
		switch c.(type) {
		case float32, float64:
			vals := s.Values.(floatValues)
			v := reflect.ValueOf(c).Float()
			return vals.addConst(v), nil
		case int, int8, int16, int32, int64:
			vals := s.Values.(floatValues)
			v := reflect.ValueOf(c).Int()
			return vals.addConst(float64(v)), nil
		default:
			return s, fmt.Errorf("Cannot add type %T to Float Series", c)
		}
	default:
		return s, fmt.Errorf("AddConst not Series  %v", s.Kind)
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
		return math.NaN(), fmt.Errorf("Mean not Series  %v", s.Kind)
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
		return math.NaN(), fmt.Errorf("Median not Series  %v", s.Kind)
	}
}

func (s Series) Min() (float64, error) {
	switch s.Kind {
	case Float:
		vals := s.Values.(floatValues)
		return vals.min(), nil
	case Int:
		vals := s.Values.(intValues)
		return vals.min(), nil
	default:
		return math.NaN(), fmt.Errorf("Min not supported for Series type %v", s.Kind)
	}
}

func (s Series) Max() (float64, error) {
	switch s.Kind {
	case Float:
		vals := s.Values.(floatValues)
		return vals.max(), nil
	case Int:
		vals := s.Values.(intValues)
		return vals.max(), nil
	default:
		return math.NaN(), fmt.Errorf("Max not supported for Series type %v", s.Kind)
	}
}

func (s Series) ValueCounts() (map[string]int, error) {
	switch s.Kind {
	case String:
		vals := s.Values.(stringValues)
		return vals.valueCounts(), nil
	default:
		return nil, fmt.Errorf("ValueCounts not supported for Series type %v", s.Kind)
	}
}

func (s Series) Count() int {
	switch s.Kind {
	case Float:
		vals := s.Values.(floatValues)
		return vals.count()
	case Int:
		vals := s.Values.(intValues)
		return vals.count()
	case Bool:
		vals := s.Values.(boolValues)
		return vals.count()
	case String:
		vals := s.Values.(stringValues)
		return vals.count()
	case DateTime:
		vals := s.Values.(dateTimeValues)
		return vals.count()
	default:
		log.Printf("Count not supported for Series type %v", s.Kind)
		return 0
	}
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
// each struct must contain a boolean field called "v"
func print(data interface{}) string {
	d := reflect.ValueOf(data)
	var printer string
	for i := 0; i < d.Len(); i++ {
		printer += fmt.Sprintln(d.Index(i).FieldByName("v"))
	}
	return printer
}
