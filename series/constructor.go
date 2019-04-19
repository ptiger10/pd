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
	name  string
}

func Type(t reflect.Kind) newSeriesOption {
	return func(c *newSeriesConfig) {
		c.kind = t
	}
}

func Name(n string) newSeriesOption {
	return func(c *newSeriesConfig) {
		c.name = n
	}
}

// New Series constructor
// that expects to receive a slice of values.
// If passing []interface{}, must supply a type expectation for the Series.
// Options: Float, Int, String, Bool, DateTime
func New(data interface{}, options ...newSeriesOption) (Series, error) {
	config := newSeriesConfig{kind: None}
	for _, option := range options {
		option(&config)
	}
	s := Series{
		Kind: config.kind,
		Name: config.name,
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
		switch config.kind {
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
			return s, fmt.Errorf("Type not supported for conversion from []interface: %v", config.kind)
		}

	default:
		return s, fmt.Errorf("Type not supported: %T", data)
	}

	s.Index =
		Index{
			Levels: []IndexLevel{
				IndexLevel{
					Type: Int,
					Labels: intLabels{
						l: makeRange(0, s.Count()),
					},
				}}}
	return s, nil
}

func makeRange(min, max int) []int64 {
	a := make([]int64, max-min)
	for i := range a {
		a[i] = int64(min + i)
	}
	return a
}
