package constructors

import (
	"fmt"
	"reflect"
	"time"

	"github.com/ptiger10/pd/new/internal/index"
	"github.com/ptiger10/pd/new/internal/values"
)

// LevelFromSlice creates an Index Level from an interface{} that reflects Slice
func LevelFromSlice(data interface{}, name string) (index.Level, error) {
	switch data.(type) {
	case []float32, []float64:
		return SliceFloat(data, name), nil
	case []int, []int8, []int16, []int32, []int64:
		return SliceInt(data, name), nil
	case []uint, []uint8, []uint16, []uint32, []uint64:
		return SliceUint(data, name), nil
	case []string:
		return SliceString(data, name), nil
	case []bool:
		return SliceBool(data, name), nil
	case []time.Time:
		return SliceDateTime(data, name), nil
	case []interface{}:
		return SliceInterface(data, name), nil
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
