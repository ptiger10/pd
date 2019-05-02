package constructors

import (
	"fmt"
	"reflect"
	"time"

	"github.com/ptiger10/pd/new/internal/index"
	"github.com/ptiger10/pd/new/internal/values"
	constructVal "github.com/ptiger10/pd/new/internal/values/constructors"
)

// LevelFromSlice creates an Index Level from an interface{} that reflects Slice
func LevelFromSlice(data interface{}, name string) (index.Level, error) {
	switch data.(type) {
	case []float32, []float64:
		vals := constructVal.InterfaceToSliceFloat(data)
		return SliceFloat(vals, name), nil
	case []int, []int8, []int16, []int32, []int64:
		vals := constructVal.InterfaceToSliceInt(data)
		return SliceInt(vals, name), nil
	case []uint, []uint8, []uint16, []uint32, []uint64:
		vals := constructVal.InterfaceToSliceUint(data)
		return SliceInt(vals, name), nil
	case []string:
		vals := data.([]string)
		return SliceString(vals, name), nil
	case []bool:
		vals := data.([]bool)
		return SliceBool(vals, name), nil
	case []time.Time:
		vals := data.([]time.Time)
		return SliceDateTime(vals, name), nil
	case []interface{}:
		vals := data.([]interface{})
		return SliceInterface(vals, name), nil
	default:
		return index.Level{}, fmt.Errorf("Unable to create level from Slice: data type not supported: %T", data)
	}
}

// Level returns an Index level with updated label map and longest value computed.
// NB: Create labels using the values.constructors methods
func level(labels values.Values, kind reflect.Kind, name string) index.Level {
	lvl := index.Level{
		Labels: labels,
		Kind:   kind,
		Name:   name,
	}
	lvl.Refresh()
	return lvl
}
