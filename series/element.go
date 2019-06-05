package series

import (
	"github.com/ptiger10/pd/kinds"
)

// An Elem is a single item in a Series
type Elem struct {
	Value      interface{}
	Null       bool
	Kind       kinds.Kind
	Index      []interface{}
	IndexKinds []kinds.Kind
}

// Elem returns the Series Element at position
func (s Series) Elem(position int) Elem {
	elem := s.values.Element(position)
	kind := s.kind

	var idx []interface{}
	var idxKinds []kinds.Kind
	for _, lvl := range s.index.Levels {
		idxElem := lvl.Labels.Element(position)
		idxVal := idxElem.Value
		idx = append(idx, idxVal)
		idxKinds = append(idxKinds, lvl.Kind)
	}
	return Elem{elem.Value, elem.Null, kind, idx, idxKinds}
}
