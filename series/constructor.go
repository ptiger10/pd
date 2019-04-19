package series

import (
	"fmt"
	"reflect"
	"time"
)

type newSeriesOption func(*newSeriesConfig)
type newSeriesConfig struct {
	kind  reflect.Kind
	index interface{}
}

func SeriesType(t reflect.Kind) newSeriesOption {
	return func(c *newSeriesConfig) {
		c.kind = t
	}
}

// New Series constructor
// that expects to receive a slice of values.
// If passing []interface{}, must supply a type expectation for the Series.
// Options: Float, Int, String, Bool, DateTime
func New(data interface{}, options ...newSeriesOption) (Series, error) {
	advanced := newSeriesConfig{kind: None}
	for _, option := range options {
		option(&advanced)
	}
	s := Series{
		Kind: advanced.kind,
	}

	switch data.(type) {
	case []float32, []float64:
		vals := floatToFloatValues(data)
		s.Values = vals
		s.Kind = Float

	case []int, []int8, []int16, []int32, []int64:
		vals := intToIntValues(data)
		s.Values = vals
		s.Kind = Int

	case []uint, []uint8, []uint16, []uint32, []uint64:
		vals := uIntToIntValues(data)
		s.Values = vals
		s.Kind = Int

	case []string:
		vals := stringToStringValues(data)
		s.Values = vals
		s.Kind = String

	case []bool:
		vals := boolToBoolValues(data)
		s.Values = vals
		s.Kind = Bool

	case []time.Time:
		vals := timeToDateTimeValues(data)
		s.Values = vals
		s.Kind = DateTime

	case []interface{}:
		d := reflect.ValueOf(data)
		switch advanced.kind {
		case None: // this checks for the pseduo-nil type
			return Series{}, fmt.Errorf("Must supply a SeriesType to decode interface")
		case Float:
			vals := interfaceToFloatValues(d)
			s.Values = vals
		case Int:
			vals := interfaceToIntValues(d)
			s.Values = vals
		case String:
			vals := interfaceToStringValues(d)
			s.Values = vals
		case Bool:
			vals := interfaceToBoolValues(d)
			s.Values = vals
		case DateTime:
			vals := interfaceToDateTimeValues(d)
			s.Values = vals
		default:
			return s, fmt.Errorf("Type not supported for conversion from []interface: %v", advanced.kind)
		}

	default:
		return s, fmt.Errorf("Type not supported: %T", data)
	}
	return s, nil
}
