package index

import (
	"fmt"

	"github.com/ptiger10/pd/internal/values"
	"github.com/ptiger10/pd/kinds"
)

// NewLevelFromSlice creates an Index Level from a Slice interface{}
// data MUST be reflect.Kind = Slice
func NewLevelFromSlice(data interface{}, name string) (Level, error) {
	vf, err := values.SliceFactory(data)
	if err != nil {
		return Level{}, fmt.Errorf("unable to create level from Slice: data type not supported: %T", data)
	}
	return newLevel(vf.V, vf.Kind, name), nil

}

// newLevel returns an Index level with updated label map and longest value computed.
// NB: Create labels using the values.constructors factory methods
func newLevel(labels values.Values, kind kinds.Kind, name string) Level {
	lvl := Level{
		Labels: labels,
		Kind:   kind,
		Name:   name,
	}
	lvl.Refresh()
	return lvl
}
