package index

import (
	"reflect"

	"github.com/ptiger10/pd/new/internal/values"
	"github.com/ptiger10/pd/new/options"
)

type Index struct {
	Levels  []Level
	NameMap LabelMap
}

type Level struct {
	Kind     reflect.Kind
	Labels   values.Values
	LabelMap LabelMap
	Name     string
	Longest  int
}

type Labels interface {
	In([]int) interface{}
}

type LabelMap map[string][]int

// ComputeLongest finds the max length of either the level name or the longest string in the LabelMap,
// for use in printing a Series or DataFrame
func (lvl *Level) ComputeLongest() {
	var max int
	for k := range lvl.LabelMap {
		if len(k) > max {
			max = len(k)
		}
	}
	if len(lvl.Name) > max {
		max = len(lvl.Name)
	}
	if max > options.DisplayIndexMaxWidth {
		max = options.DisplayIndexMaxWidth
	}
	lvl.Longest = max
}
