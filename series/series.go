package series

import (
	"fmt"
	"reflect"

	"github.com/ptiger10/pd/internal/index"
	"github.com/ptiger10/pd/internal/values"
	"github.com/ptiger10/pd/kinds"
)

// A Series is a 1-D data container with a labeled index, static type, and the ability to handle null values
type Series struct {
	index  index.Index
	values values.Values
	kind   kinds.Kind
	Name   string
	Math
	To
	IndexTo
}

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

// Kind is the data kind of the Series' values. Mimics reflect.Kind with the addition of time.Time
func (s Series) Kind() string {
	return fmt.Sprint(s.kind)
}

func (s Series) copy() Series {
	idx := s.index.Copy()
	valsCopy := s.values.Copy()
	copyS := &Series{
		values: valsCopy,
		index:  idx,
		kind:   s.kind,
		Name:   s.Name,
	}
	copyS.Math = Math{s: copyS}
	copyS.To = To{s: copyS}
	copyS.IndexTo = IndexTo{s: copyS}
	return *copyS
}

// in subsets a Series to include only index items and values at the positions supplied,
// then returns as a new Series
func (s Series) in(positions []int) (Series, error) {
	if ok := s.ensureAlignment(); !ok {
		return s, fmt.Errorf("fatal error: Series values and index labels out of alignment: report issue and create new series")
	}
	newS := s.copy()
	values, err := newS.values.In(positions)
	if err != nil {
		return Series{}, fmt.Errorf("unable to get Series values at positions: %v", err)
	}
	newS.values = values
	for i, level := range newS.index.Levels {
		// Ducks error because positional alignment is ensured between values and all index levels
		newS.index.Levels[i].Labels, _ = level.Labels.In(positions)
	}
	newS.index.Refresh()
	return newS, nil
}

func seriesEquals(s1, s2 Series) bool {
	sameIndex := reflect.DeepEqual(s1.index, s2.index)
	sameValues := reflect.DeepEqual(s1.values, s2.values)
	sameName := s1.Name == s2.Name
	sameKind := s1.kind == s2.kind
	if sameIndex && sameValues && sameName && sameKind {
		return true
	}
	return false
}
