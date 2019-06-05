package series

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/ptiger10/pd/internal/index"
	"github.com/ptiger10/pd/internal/values"
	"github.com/ptiger10/pd/kinds"
	"github.com/ptiger10/pd/opt"
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

// An Element is a single item in a Series
type Element struct {
	Value      interface{}
	Null       bool
	Kind       kinds.Kind
	Labels     []interface{}
	LabelKinds []kinds.Kind
}

func (el Element) String() string {
	var printStr string
	for _, pair := range [][]interface{}{
		[]interface{}{"Value", el.Value},
		[]interface{}{"Null", el.Null},
		[]interface{}{"Kind", el.Kind},
		[]interface{}{"Labels", el.Labels},
		[]interface{}{"LabelKinds", el.LabelKinds},
	} {
		// LabelKinds is 10 characters wide, so left padding set to 10
		printStr += fmt.Sprintf("%10v:%v%v\n", pair[0], strings.Repeat(" ", opt.GetDisplayElementWhitespaceBuffer()), pair[1])
	}
	return printStr
}

// Element returns the Series Element at position
func (s Series) Element(position int) Element {
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
	return Element{elem.Value, elem.Null, kind, idx, idxKinds}
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

// in copies a Series then subsets it to include only index items and values at the positions supplied
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
