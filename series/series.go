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
	Type   reflect.Kind
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

type Int []int64

func (vals floatValues) Sum() (float64, error) {
	var sum float64
	for _, val := range vals {
		sum += val.v
	}
	return sum, nil
}

func (vals Int) Sum() (float64, error) {
	var sum float64
	for _, val := range vals {
		sum += float64(val)
	}
	return sum, nil
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

const Float = reflect.Float64

func New(data interface{}, options ...newSeriesOption) (Series, error) {
	advanced := newSeriesConfig{kind: reflect.UnsafePointer}
	for _, option := range options {
		option(&advanced)
	}
	s := Series{}

	d := reflect.ValueOf(data)
	switch data.(type) {
	case []float32, []float64:
		var vals []floatValue
		for i := 0; i < d.Len(); i++ {
			vals = append(vals, floatValue{v: d.Index(i).Float()})
			s.Values = floatValues(vals)
		}
	case []int, []int8, []int16, []int32, []int64:
		var vals Int
		for i := 0; i < d.Len(); i++ {
			vals = append(vals, d.Index(i).Int())
			s.Values = vals
		}

	case []uint, []uint8, []uint16, []uint32, []uint64:
		var vals Int
		for i := 0; i < d.Len(); i++ {
			vals = append(vals, int64(d.Index(i).Uint()))
			s.Values = vals
		}
	case []interface{}:
		switch advanced.kind {
		case reflect.UnsafePointer:
			return Series{}, fmt.Errorf("Must supply a SeriesType to decode interface")
		case Float:
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
				}
			}
			s.Values = floatValues(vals)
		}

	default:
		return s, fmt.Errorf("Data type not supported: %T", data)
	}
	return s, nil
}
