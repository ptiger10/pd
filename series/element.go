package series

import (
	"github.com/ptiger10/pd/kinds"
)

// An Element is a single item in a Series
type Element struct {
	Value      interface{}
	Null       bool
	Kind       kinds.Kind
	Index      []interface{}
	IndexKinds []kinds.Kind
}

// Elem returns the Series Element at position
func (s Series) Elem(position int) Element {
	elemSlice := s.values.Element(position)
	val := elemSlice[0]
	null := elemSlice[1].(bool)
	kind := s.Kind

	var idx []interface{}
	var idxKinds []kinds.Kind
	for _, lvl := range s.index.Levels {
		idxElem := lvl.Labels.Element(position)
		idxVal := idxElem[0]
		idx = append(idx, idxVal)
		idxKinds = append(idxKinds, lvl.Kind)
	}
	return Element{val, null, kind, idx, idxKinds}
}
