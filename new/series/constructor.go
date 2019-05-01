package series

import (
	"fmt"
	"log"
	"reflect"

	"github.com/ptiger10/pd/new/internal/index"
	constructIdx "github.com/ptiger10/pd/new/internal/index/constructors"
	"github.com/ptiger10/pd/new/internal/values"
	constructVal "github.com/ptiger10/pd/new/internal/values/constructors"
	"github.com/ptiger10/pd/new/kinds"
)

type Option func(*seriesConfig)
type seriesConfig struct {
	kind  reflect.Kind
	index []idx
	name  string
}

type idx struct {
	data interface{}
	name string
}

func Kind(t reflect.Kind) Option {
	return func(c *seriesConfig) {
		c.kind = t
	}
}

func Name(n string) Option {
	return func(c *seriesConfig) {
		c.name = n
	}
}

// Index returns a Option for use in the Series constructor New(),
// and takes an optional Name. If name is blank, defaults to RangeLabels (0, 1, 2, ...n)
func Index(data interface{}, options ...Option) Option {
	config := seriesConfig{}
	for _, option := range options {
		option(&config)
	}
	return func(c *seriesConfig) {
		idx := idx{
			data: data,
			name: config.name,
		}
		c.index = append(c.index, idx)

	}
}

// Calls New and panics if error. For use in testing
func mustNew(data interface{}, options ...Option) Series {
	s, err := New(data, options...)
	if err != nil {
		log.Panic(err)
	}
	return s
}

// New Series constructor
// Optional
// - Index(): If no index is supplied, defaults to RangeLabels (0, 1, 2, ...n)
// - Name(): If no name is supplied, no name will appear when Series is printed
// - Kind(): Supplying a type will try to cast the Series values into a specific kind
// If passing []interface{}, must supply a type expectation for the Series.
// Options: Float, Int, String, Bool, DateTime
func New(data interface{}, options ...Option) (Series, error) {
	// config := seriesConfig{kind: datatypes.None}
	config := seriesConfig{}

	for _, option := range options {
		option(&config)
	}
	var v values.Values
	var idx index.Index
	kind := config.kind
	name := config.name

	switch data.(type) {
	// case []float32, []float64:
	// 	vals := floatToFloatValues(data)
	// 	s.Values = vals
	// 	s.Kind = Float

	case []int, []int8, []int16, []int32, []int64:
		v = constructVal.SliceInt(data)
		kind = kinds.Int

	case []uint, []uint8, []uint16, []uint32, []uint64:
		v = constructVal.SliceUInt(data)
		kind = kinds.Int

	case []string:
		v = constructVal.SliceString(data)
		kind = kinds.String

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
		return Series{}, fmt.Errorf("Type not supported: %T", data)
	}

	if config.index == nil {
		n := len(v.All())
		idx = constructIdx.Default(n)
	} else {
		var levels []index.Level
		for _, labels := range config.index {
			switch labels.data.(type) {
			case []int, []int8, []int16, []int32, []int64:
				level := constructIdx.SliceInt(labels.data, labels.name)
				levels = append(levels, level)
			}
		}
		idx = constructIdx.New(levels...)
	}
	s := Series{
		Index:  idx,
		Values: v,
		Kind:   kind,
		Name:   name,
	}

	return s, nil
}

func (s Series) Set(subset Series, data interface{}) {

}
