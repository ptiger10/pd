package constructors

import (
	"fmt"
	"reflect"
	"time"

	"github.com/ptiger10/pd/new/internal/values"
	"github.com/ptiger10/pd/new/kinds"
)

// ValuesFromSlice creates values from an interface{} with Slice reflection
func ValuesFromSlice(data interface{}) (values.Values, reflect.Kind, error) {
	var v values.Values
	var kind reflect.Kind

	switch data.(type) {
	case []float32, []float64:
		var vals []float64
		d := reflect.ValueOf(data)
		for i := 0; i < d.Len(); i++ {
			vals = append(vals, d.Index(i).Float())
		}
		v = SliceFloat(vals)
		kind = kinds.Float

	case []int, []int8, []int16, []int32, []int64:
		var vals []int64
		d := reflect.ValueOf(data)
		for i := 0; i < d.Len(); i++ {
			vals = append(vals, d.Index(i).Int())
		}
		v = SliceInt(vals)
		kind = kinds.Int

	case []uint, []uint8, []uint16, []uint32, []uint64:
		var vals []int64
		d := reflect.ValueOf(data)
		for i := 0; i < d.Len(); i++ {
			vals = append(vals, int64(d.Index(i).Uint()))
		}
		v = SliceInt(data)
		kind = kinds.Int

	case []string:
		v = SliceString(data)
		kind = kinds.String

	case []bool:
		v = SliceBool(data)
		kind = kinds.Bool

	case []time.Time:
		v = SliceDateTime(data)
		kind = kinds.DateTime

	case []interface{}:
		v = SliceInterface(data)
		kind = kinds.Interface

	default:
		return nil, kinds.None, fmt.Errorf("Type %T not supported", data)
	}

	return v, kind, nil
}
