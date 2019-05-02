package series

import "reflect"

// An Element is a single item in a Series
type Element struct {
	Value      interface{}
	Null       bool
	Kind       reflect.Kind
	Index      []interface{}
	IndexKinds []reflect.Kind
}

// Elem returns the Series Element at position
func (s Series) Elem(position int) Element {
	elemSlice := s.Values.Element(position)
	val := elemSlice[0]
	null := elemSlice[1].(bool)
	kind := s.Kind

	var idx []interface{}
	var idxKinds []reflect.Kind
	for _, lvl := range s.Index.Levels {
		idxElem := lvl.Labels.Element(position)
		idxVal := idxElem[0]
		idx = append(idx, idxVal)
		idxKinds = append(idxKinds, lvl.Kind)

	}

	return Element{val, null, kind, idx, idxKinds}
}
