package index

import (
	"fmt"
	"log"

	"github.com/ptiger10/pd/internal/values"
	"github.com/ptiger10/pd/options"
)

// A Level is a single collection of labels within an index, plus label mappings and metadata
type Level struct {
	DataType options.DataType
	Labels   values.Values
	LabelMap LabelMap
	Name     string
}

// A LabelMap records the position of labels, in the form {label name: [label position(s)]}
type LabelMap map[string][]int

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
func newLevel(labels values.Values, datatype options.DataType, name string) Level {
	lvl := Level{
		Labels:   labels,
		DataType: datatype,
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

type Element struct {
	Label    interface{}
	DataType options.DataType
}

func (lvl Level) Element(position int) Element {
	return Element{
		Label:    lvl.Labels.Element(position).Value,
		DataType: lvl.DataType,
	}
}

// MaxWidth finds the max length of either the level name or the longest string in the LabelMap,
// for use in printing a Series or DataFrame
func (lvl *Level) MaxWidth() int {
	var max int
	for k := range lvl.LabelMap {
		if len(k) > max {
			max = len(k)
		}
	}
	if len(lvl.Name) > max {
		max = len(lvl.Name)
	}
	if max > options.GetDisplayMaxWidth() {
		max = options.GetDisplayMaxWidth()
	}
	return max
}
