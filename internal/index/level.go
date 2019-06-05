package index

import (
	"fmt"
	"reflect"

	"github.com/ptiger10/pd/internal/values"
	"github.com/ptiger10/pd/kinds"
)

// NewLevel creates an Index Level from a Scalar or Slice interface{}
func NewLevel(data interface{}, name string) (Level, error) {
	var vf values.Factory
	var err error

	switch reflect.ValueOf(data).Kind() {
	case reflect.Slice:
		vf, err = values.SliceFactory(data)
		if err != nil {
			return Level{}, fmt.Errorf("unable to create level from Slice: %v", data)
		}
	default:
		vf, err = values.ScalarFactory(data)
		if err != nil {
			return Level{}, fmt.Errorf("unable to create level from Scalar: %v", data)
		}
	}
	return newLevel(vf.Values, vf.Kind, name), nil
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

// Copy copies an Index Level
func (lvl Level) Copy() Level {
	lvlCopy := Level{}
	lvlCopy = lvl
	lvlCopy.Labels = lvlCopy.Labels.Copy()
	lvlCopy.LabelMap = make(LabelMap)
	for k, v := range lvlCopy.LabelMap {
		lvlCopy.LabelMap[k] = v
	}
	return lvlCopy
}
