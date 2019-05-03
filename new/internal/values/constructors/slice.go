package constructors

import (
	"fmt"
	"reflect"
	"time"

	"github.com/ptiger10/pd/new/internal/values"
	"github.com/ptiger10/pd/new/kinds"
)

// ExtendedValues is an extended representation of values containing values and kind
type ExtendedValues struct {
	V    values.Values
	Kind kinds.Kind
}

// ValuesFromSlice creates values from an interface{} with Slice reflection
func ValuesFromSlice(data interface{}) (ExtendedValues, error) {
	var v values.Values
	var kind kinds.Kind

	switch data.(type) {
	case []float32, []float64:
		vals := InterfaceToSliceFloat(data)
		v = SliceFloat(vals)
		kind = kinds.Float

	case []int, []int8, []int16, []int32, []int64:
		vals := InterfaceToSliceInt(data)
		v = SliceInt(vals)
		kind = kinds.Int

	case []uint, []uint8, []uint16, []uint32, []uint64:
		// returns []int64 so it can be passed into SliceInt
		vals := InterfaceToSliceUint(data)
		v = SliceInt(vals)
		kind = kinds.Int

	case []string:
		vals := data.([]string)
		v = SliceString(vals)
		kind = kinds.String

	case []bool:
		vals := data.([]bool)
		v = SliceBool(vals)
		kind = kinds.Bool

	case []time.Time:
		vals := data.([]time.Time)
		v = SliceDateTime(vals)
		kind = kinds.DateTime

	case []interface{}:
		vals := data.([]interface{})
		v = SliceInterface(vals)
		kind = kinds.Interface

	default:
		ret := ExtendedValues{nil, kinds.Invalid}
		return ret, fmt.Errorf("Type %T not supported", data)
	}

	ret := ExtendedValues{v, kind}
	return ret, nil
}

// [START interface converters]

// InterfaceToSliceFloat converts interface{} -> []float64
func InterfaceToSliceFloat(data interface{}) []float64 {
	var vals []float64
	d := reflect.ValueOf(data)
	for i := 0; i < d.Len(); i++ {
		vals = append(vals, d.Index(i).Float())
	}
	return vals
}

// InterfaceToSliceInt converts interface{} -> []int64
func InterfaceToSliceInt(data interface{}) []int64 {
	var vals []int64
	d := reflect.ValueOf(data)
	for i := 0; i < d.Len(); i++ {
		vals = append(vals, int64(d.Index(i).Int()))
	}
	return vals
}

// InterfaceToSliceUint converts interface{} -> []int64
func InterfaceToSliceUint(data interface{}) []int64 {
	var vals []int64
	d := reflect.ValueOf(data)
	for i := 0; i < d.Len(); i++ {
		vals = append(vals, int64(d.Index(i).Uint()))
	}
	return vals
}

// [END interface converters]
