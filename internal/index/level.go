package index

import (
	"fmt"
	"log"

	"github.com/ptiger10/pd/internal/values"
	"github.com/ptiger10/pd/options"
)

// NewLevel creates an Index Level from a Scalar or Slice interface{} but returns an error if interface{} is not supported by factory.
func NewLevel(data interface{}, name string) (Level, error) {
	factory, err := values.InterfaceFactory(data)
	if err != nil {
		return Level{}, fmt.Errorf("NewLevel(): %v", err)
	}
	return newLevel(factory.Values, factory.DataType, name), nil
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
func newLevel(labels values.Values, kind options.DataType, name string) Level {
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
	for k, v := range lvl.LabelMap {
		lvlCopy.LabelMap[k] = v
	}
	return lvlCopy
}

// Len returns the number of labels in the level
func (lvl Level) Len() int {
	return lvl.Labels.Len()
}

// DefaultLevel creates an unnamed index level with range labels (0, 1, 2, ...n)
func DefaultLevel(n int, name string) Level {
	v := values.NewDefault(n)
	level := newLevel(v, options.Int64, name)
	return level
}
