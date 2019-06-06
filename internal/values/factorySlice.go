package values

import (
	"fmt"
	"reflect"
	"time"

	"github.com/ptiger10/pd/kinds"
)

// SliceFactory creates a Factory from an interface{} that has been determined elsewhere to be a Slice.
// Be careful modifying this constructor and its children,
// as reflection is inherently unsafe and each function expects to receive specific types only.
func SliceFactory(data interface{}) (Factory, error) {
	var ret Factory

	switch data.(type) {
	case []float32:
		vals := sliceFloatToSliceFloat64(data)
		ret = newSliceFloat64(vals)

	case []float64:
		vals := data.([]float64)
		ret = newSliceFloat64(vals)

	case []int, []int8, []int16, []int32:
		vals := sliceIntToSliceInt64(data)
		ret = newSliceInt64(vals)
	case []int64:
		vals := data.([]int64)
		ret = newSliceInt64(vals)

	case []uint, []uint8, []uint16, []uint32, []uint64:
		// converts into []int64 so it can be passed into SliceInt
		vals := sliceUIntToSliceInt64(data)
		ret = newSliceInt64(vals)

	case []string:
		vals := data.([]string)
		ret = newSliceString(vals)

	case []bool:
		vals := data.([]bool)
		ret = newSliceBool(vals)

	case []time.Time:
		vals := data.([]time.Time)
		ret = newSliceDateTime(vals)

	case []interface{}:
		vals := data.([]interface{})
		ret = newSliceInterface(vals)

	default:
		ret = Factory{nil, kinds.None}
		return ret, fmt.Errorf("Type %T not supported", data)
	}

	return ret, nil
}

// [START interface converters]

// sliceFloatToSliceFloat64 converts known []float interface{} -> []float64
func sliceFloatToSliceFloat64(data interface{}) []float64 {
	var vals []float64
	d := reflect.ValueOf(data)
	for i := 0; i < d.Len(); i++ {
		v := d.Index(i).Float()
		vals = append(vals, v)
	}
	return vals
}

// sliceIntToSliceInt64 converts known []int interface{} -> []int64
func sliceIntToSliceInt64(data interface{}) []int64 {
	var vals []int64
	d := reflect.ValueOf(data)
	for i := 0; i < d.Len(); i++ {
		v := d.Index(i).Int()
		vals = append(vals, v)
	}
	return vals
}

// sliceUIntToSliceInt64 converts knonw []uint interface{} -> []int64
func sliceUIntToSliceInt64(data interface{}) []int64 {
	var vals []int64
	d := reflect.ValueOf(data)
	for i := 0; i < d.Len(); i++ {
		v := d.Index(i).Uint()
		vals = append(vals, int64(v))
	}
	return vals
}

// [END interface converters]
