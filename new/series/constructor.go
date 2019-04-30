package series

import (
	"fmt"
	"reflect"

	"github.com/ptiger10/pd/new/datatypes"
	"github.com/ptiger10/pd/new/internal/index"
	constructIdx "github.com/ptiger10/pd/new/internal/index/constructors"
	constructVal "github.com/ptiger10/pd/new/internal/values/constructors"
)

type seriesOption func(*seriesConfig)
type seriesConfig struct {
	kind  reflect.Kind
	index []interface{}
	name  string
}

func Type(t reflect.Kind) seriesOption {
	return func(c *seriesConfig) {
		c.kind = t
	}
}

func Name(n string) seriesOption {
	return func(c *seriesConfig) {
		c.name = n
	}
}

func Index(data interface{}) seriesOption {
	return func(c *seriesConfig) {
		c.index = append(c.index, data)
	}
}

// New Series constructor
// that expects to receive a slice of values.
// If passing []interface{}, must supply a type expectation for the Series.
// Options: Float, Int, String, Bool, DateTime
func New(data interface{}, options ...seriesOption) (Series, error) {
	config := seriesConfig{kind: datatypes.None}

	for _, option := range options {
		option(&config)
	}
	s := Series{
		Kind: config.kind,
		Name: config.name,
	}

	switch data.(type) {
	// case []float32, []float64:
	// 	vals := floatToFloatValues(data)
	// 	s.Values = vals
	// 	s.Kind = Float

	// case []int, []int8, []int16, []int32, []int64:
	// 	vals := intToIntValues(data)
	// 	s.Values = vals
	// 	s.Kind = Int

	// case []uint, []uint8, []uint16, []uint32, []uint64:
	// 	vals := uIntToIntValues(data)
	// 	s.Values = vals
	// 	s.Kind = Int

	case []string:
		vals := constructVal.SliceString(data)
		s.Values = vals

		s.Kind = datatypes.String

	// case []bool:
	// 	vals := boolToBoolValues(data)
	// 	s.Values = vals
	// 	s.Kind = Bool

	// case []time.Time:
	// 	vals := timeToDateTimeValues(data)
	// 	s.Values = vals
	// 	s.Kind = DateTime

	// case []interface{}:
	// 	d := reflect.ValueOf(data)
	// 	switch config.kind {
	// 	case None: // this checks for the pseduo-nil type
	// 		return Series{}, fmt.Errorf("Must supply a SeriesType to decode interface")
	// 	case Float:
	// 		vals := interfaceToFloatValues(d)
	// 		s.Values = vals
	// 	case Int:
	// 		vals := interfaceToIntValues(d)
	// 		s.Values = vals
	// 	case String:
	// 		vals := interfaceToStringValues(d)
	// 		s.Values = vals
	// 	case Bool:
	// 		vals := interfaceToBoolValues(d)
	// 		s.Values = vals
	// 	case DateTime:
	// 		vals := interfaceToDateTimeValues(d)
	// 		s.Values = vals
	// 	default:
	// 		return s, fmt.Errorf("Type not supported for conversion from []interface: %v", config.kind)
	// 	}

	default:
		return s, fmt.Errorf("Type not supported: %T", data)
	}

	if config.index == nil {
		level := constructIdx.SliceInt(makeRange(0, s.Len()))
		s.Index = constructIdx.New([]index.Level{level})
	} else {
		var levels []index.Level
		for _, idxInput := range config.index {
			switch idxInput.(type) {
			case []int, []int8, []int16, []int32, []int64:
				level := constructIdx.SliceInt(idxInput)
				levels = append(levels, level)
			}
		}
		s.Index = constructIdx.New(levels)
	}

	return s, nil
}

func (s Series) Set(subset Series, data interface{}) {

}

func makeRange(min, max int) []int64 {
	a := make([]int64, max-min)
	for i := range a {
		a[i] = int64(min + i)
	}
	return a
}
