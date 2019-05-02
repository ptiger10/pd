package index

import (
	"reflect"

	"github.com/ptiger10/pd/new/internal/values"
)

// An Index is a collection of levels, plus label mappings
type Index struct {
	Levels  []Level
	NameMap LabelMap
}

// A Level is a single collection of labels within an index, plus label mappings and metadata
type Level struct {
	Kind     reflect.Kind
	Labels   values.Values
	LabelMap LabelMap
	Name     string
	Longest  int
}

// A LabelMap records the position of labels, in the form {label name: [label position(s)]}
type LabelMap map[string][]int

// A MiniIndex is a mini representation of one index level.
// It is used for unpacking client-supplied index data and optional metadata
type MiniIndex struct {
	Data interface{}
	Kind reflect.Kind
	Name string
}
