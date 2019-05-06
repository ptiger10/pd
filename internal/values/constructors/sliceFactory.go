package constructors

import (
	"fmt"
	"reflect"
	"time"

	"github.com/ptiger10/pd/internal/values"
	"github.com/ptiger10/pd/kinds"
)

// ValuesFactory is an extended representation of values containing values and kind
type ValuesFactory struct {
	V    values.Values
	Kind kinds.Kind
}

// ValuesFromSlice creates values from an interface{} that has been determined elsewhere to be a Slice.
// Be careful modifying this constructor and its children,
// as reflection is inherently unsafe and each function expects to receive specific types only.
func ValuesFromSlice(data interface{}) (ValuesFactory, error) {
	var ret ValuesFactory

	switch data.(type) {
	case []float32:
		vals := InterfaceToSliceFloat(data)
		ret = SliceFloat(vals)

	case []float64:
		vals := data.([]float64)
		ret = SliceFloat(vals)

	case []int, []int8, []int16, []int32:
		vals := InterfaceToSliceInt(data)
		ret = SliceInt(vals)

	case []int64:
		vals := data.([]int64)
		ret = SliceInt(vals)

	case []uint, []uint8, []uint16, []uint32, []uint64:
		// returns []int64 so it can be passed into SliceInt
		vals := InterfaceToSliceUint(data)
		ret = SliceInt(vals)

	case []string:
		vals := data.([]string)
		ret = SliceString(vals)

	case []bool:
		vals := data.([]bool)
		ret = SliceBool(vals)

	case []time.Time:
		vals := data.([]time.Time)
		ret = SliceDateTime(vals)

	case []interface{}:
		vals := data.([]interface{})
		ret = SliceInterface(vals)

	default:
		ret = ValuesFactory{nil, kinds.Invalid}
		return ret, fmt.Errorf("Type %T not supported", data)
	}

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
