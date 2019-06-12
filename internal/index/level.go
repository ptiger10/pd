package index

import (
	"fmt"
	"log"
	"reflect"

	"github.com/ptiger10/pd/datatypes"
	"github.com/ptiger10/pd/internal/values"
)

// NewLevel creates an Index Level from a Scalar or Slice interface{} but returns an error if interface{} is not supported by factory.
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
	return newLevel(vf.Values, vf.DataType, name), nil
}

// MustCreateNewLevel returns a new level from an interface, but panics on error
func MustCreateNewLevel(data interface{}, name string) Level {
	lvl, err := NewLevel(data, name)
	if err != nil {
		log.Fatalf("MustCreateNewLevel returned an error: %v", err)
	}
	return lvl
}

// newLevel returns an Index level with updated label map and longest value computed. Never returns an error.
// NB: Create labels using the values.constructors factory methods, as in NewLevel().
func newLevel(labels values.Values, kind datatypes.DataType, name string) Level {
	lvl := Level{
		Labels:   labels,
		DataType: kind,
		Name:     name,
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

// Len returns the number of labels in the level
func (lvl Level) Len() int {
	return lvl.Labels.Len()
}
