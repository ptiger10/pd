package index

import (
	"github.com/ptiger10/pd/internal/values"
	"github.com/ptiger10/pd/kinds"
)

// An Index is a collection of levels, plus label mappings
type Index struct {
	Levels  []Level
	NameMap LabelMap
}

// A Level is a single collection of labels within an index, plus label mappings and metadata
type Level struct {
	Kind     kinds.Kind
	Labels   values.Values
	LabelMap LabelMap
	Name     string
	Longest  int
}

// A LabelMap records the position of labels, in the form {label name: [label position(s)]}
type LabelMap map[string][]int
