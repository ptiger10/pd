package series

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
)

type Series struct {
	Index  Index
	Values Values
	Kind   reflect.Kind
}

func (s Series) Sum() (float64, error) {
	return s.Values.Sum()
}

type Index struct {
	Levels []IndexLevel
}

type IndexLevel struct {
	Type   reflect.Kind
	Values Values
}

type Values interface {
	Sum() (float64, error)
}

type floatValues []floatValue
type floatValue struct {
	v    float64
	null bool
}

type intValues []intValue
type intValue struct {
	v    int64
	null bool
}

type stringValues []stringValue
type stringValue struct {
	v    string
	null bool
}

func (vals floatValues) Sum() (float64, error) {
	var sum float64
	for _, val := range vals {
		sum += val.v
	}
	return sum, nil
}

func (vals intValues) Sum() (float64, error) {
	var sum float64
	for _, val := range vals {
		sum += float64(val.v)
	}
	return sum, nil
}

func (vals stringValues) Sum() (float64, error) {
	return math.NaN(), fmt.Errorf("Unable to call Sum() on Series of strings")
}

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

const (
	Float  = reflect.Float64
	Int    = reflect.Int64
	String = reflect.String
)

// pseudo-nil value for type reflect.Kind
const None = reflect.UnsafePointer

func New(data interface{}, options ...newSeriesOption) (Series, error) {
	advanced := newSeriesConfig{kind: None}
	for _, option := range options {
		option(&advanced)
	}
	s := Series{
		Kind: advanced.kind,
	}

	d := reflect.ValueOf(data)
	switch data.(type) {
	case []float32, []float64:
		vals := floatToFloatValues(d)
		s.Values = vals
		s.Kind = Float

	case []int, []int8, []int16, []int32, []int64:
		vals := intToIntValues(d)
		s.Values = vals
		s.Kind = Int

	case []uint, []uint8, []uint16, []uint32, []uint64:
		vals := uIntToIntValues(d)
		s.Values = vals
		s.Kind = Int

	case []string:
		vals := stringToStringValues(d)
		s.Values = vals
		s.Kind = String

	case []interface{}:
		switch advanced.kind {
		case None: // this checks for the pseduo-nil type
			return Series{}, fmt.Errorf("Must supply a SeriesType to decode interface")
		case Float:
			vals := interfaceToFloatValues(d)
			s.Values = vals
		}

	default:
		return s, fmt.Errorf("Data type not supported: %T", data)
	}
	return s, nil
}

func floatToFloatValues(d reflect.Value) floatValues {
	var vals []floatValue
	for i := 0; i < d.Len(); i++ {
		vals = append(vals, floatValue{v: d.Index(i).Float()})
	}
	return floatValues(vals)
}

func intToIntValues(d reflect.Value) intValues {
	var vals []intValue
	for i := 0; i < d.Len(); i++ {
		vals = append(vals, intValue{v: d.Index(i).Int()})
	}
	return intValues(vals)
}

func uIntToIntValues(d reflect.Value) intValues {
	var vals []intValue
	for i := 0; i < d.Len(); i++ {
		val := int64(d.Index(i).Uint())
		vals = append(vals, intValue{v: val})
	}
	return intValues(vals)
}

func stringToStringValues(d reflect.Value) stringValues {
	var vals []stringValue
	for i := 0; i < d.Len(); i++ {
		vals = append(vals, stringValue{v: d.Index(i).String()})
	}
	return stringValues(vals)
}

func interfaceToFloatValues(d reflect.Value) floatValues {
	var vals []floatValue
	for i := 0; i < d.Len(); i++ {
		v := d.Index(i).Elem()
		fmt.Println(v)
		switch d.Index(i).Elem().Kind() {
		case reflect.Invalid:
			vals = append(vals, floatValue{null: true})
		case reflect.String:
			val, err := strconv.ParseFloat(v.String(), 64)
			if err != nil || math.IsNaN(val) {
				vals = append(vals, floatValue{null: true})
			} else {
				vals = append(vals, floatValue{v: val})
			}
		case reflect.Float32, reflect.Float64:
			vals = append(vals, floatValue{v: v.Float()})
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			vals = append(vals, floatValue{v: float64(v.Int())})
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			vals = append(vals, floatValue{v: float64(v.Uint())})
		}
	}
	return floatValues(vals)
}
