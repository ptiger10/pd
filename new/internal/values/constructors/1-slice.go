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
		v = SliceFloat(data)
		kind = kinds.Float

	case []int, []int8, []int16, []int32, []int64:
		v = SliceInt(data)
		kind = kinds.Int

	case []uint, []uint8, []uint16, []uint32, []uint64:
		v = SliceUInt(data)
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
		return nil, kinds.None, fmt.Errorf("Datatype %T not supported", data)
	}

	return v, kind, nil
}
