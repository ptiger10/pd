package series

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/ptiger10/pd/internal/index"
	"github.com/ptiger10/pd/internal/values"
	"github.com/ptiger10/pd/kinds"
	"github.com/ptiger10/pd/opt"
)

// A Series is a 1-D data container with a labeled index, static type, and the ability to handle null values
type Series struct {
	index   index.Index
	values  values.Values
	kind    kinds.Kind
	Name    string
	Math    Math
	To      To
	Index   Index
	Select  Select
	InPlace InPlace
}

// An Element is a single item in a Series.
type Element struct {
	Value      interface{}
	Null       bool
	Labels     []interface{}
	LabelKinds []kinds.Kind
}

func (el Element) String() string {
	var printStr string
	for _, pair := range [][]interface{}{
		[]interface{}{"Value", el.Value},
		[]interface{}{"Null", el.Null},
		[]interface{}{"Labels", el.Labels},
		[]interface{}{"LabelKinds", el.LabelKinds},
	} {
		// LabelKinds is 10 characters wide, so left padding set to 10
		printStr += fmt.Sprintf("%10v:%v%v\n", pair[0], strings.Repeat(" ", opt.GetDisplayElementWhitespaceBuffer()), pair[1])
	}
	return printStr
}

// Element returns information about the value and index labels at this position.
func (s Series) Element(position int) Element {
	elem := s.values.Element(position)

	var idx []interface{}
	var idxKinds []kinds.Kind
	for _, lvl := range s.index.Levels {
		idxElem := lvl.Labels.Element(position)
		idxVal := idxElem.Value
		idx = append(idx, idxVal)
		idxKinds = append(idxKinds, lvl.Kind)
	}
	return Element{elem.Value, elem.Null, idx, idxKinds}
}

// Kind is the data kind of the Series' values. Mimics reflect.Kind with the addition of time.Time as DateTime.
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
	copyS.Index = Index{s: copyS, To: To{s: copyS, idx: true}}
	copyS.Select = Select{s: copyS}
	copyS.InPlace = InPlace{s: copyS}
	return *copyS
}

// in copies a Series then subsets it to include only index items and values at the positions supplied
func (s Series) in(positions []int) (Series, error) {
	if err := s.ensureAlignment(); err != nil {
		return s, fmt.Errorf("Series.in(): %v", err)
	}
	newS := s.copy()
	values, err := newS.values.In(positions)
	if err != nil {
		return Series{}, fmt.Errorf("Series.in() values: %v", err)
	}
	newS.values = values
	for i, level := range newS.index.Levels {
		newS.index.Levels[i].Labels, err = level.Labels.In(positions)
		if err != nil {
			return Series{}, fmt.Errorf("Series.in() index: %v", err)
		}
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

// Len returns the number of Elements (i.e., Value/Null pairs) in the Series.
func (s Series) Len() int {
	return s.values.Len()
}

// valid returns integer positions of valid (i.e., non-null) values in the series.
func (s Series) valid() []int {
	var ret []int
	for i := 0; i < s.Len(); i++ {
		if !s.values.Element(i).Null {
			ret = append(ret, i)
		}
	}
	return ret
}

// null returns the integer position of all null values in the collection.
func (s Series) null() []int {
	var ret []int
	for i := 0; i < s.Len(); i++ {
		if s.values.Element(i).Null {
			ret = append(ret, i)
		}
	}
	return ret
}

// all returns only the Value fields for the collection of Value/Null structs as an interface slice.
//
// Caution: This operation excludes the Null field but retains any null values.
func (s Series) all() []interface{} {
	var ret []interface{}
	for i := 0; i < s.Len(); i++ {
		ret = append(ret, s.values.Element(i).Value)
	}
	return ret
}

// Insert inserts a new row into the Series immediately before the specified integer position and returns a new Series.
func (s Series) Insert(pos int, val interface{}, idx []interface{}) (Series, error) {
	if len(idx) != s.index.Len() {
		return Series{}, fmt.Errorf("Series.Insert() len(idx) must equal number of index levels: supplied %v want %v",
			len(idx), s.index.Len())
	}
	s = s.copy()
	for i := 0; i < s.index.Len(); i++ {
		err := s.index.Levels[i].Labels.Insert(pos, idx[i])
		if err != nil {
			return Series{}, fmt.Errorf("Series.Insert() with idx val %v at idx level %v: %v", val, i, err)
		}
		s.index.Levels[i].Refresh()
	}
	if err := s.values.Insert(pos, val); err != nil {
		return Series{}, fmt.Errorf("Series.Insert() with val %v: %v", val, err)
	}
	return s, nil
}

// dropRows drops multiple rows and returns a new Series
func (s Series) dropRows(positions []int) (Series, error) {
	var err error
	sort.IntSlice(positions).Sort()
	for i, position := range positions {
		s, err = s.Drop(position - i)
		if err != nil {
			return Series{}, err
		}
	}
	return s, nil
}

// Drop drops the row at the specified integer position and returns a new Series.
func (s Series) Drop(pos int) (Series, error) {
	s = s.copy()
	for i := 0; i < s.index.Len(); i++ {
		err := s.index.Levels[i].Labels.Drop(pos)
		if err != nil {
			return Series{}, fmt.Errorf("Series.Drop(): %v", err)
		}
		s.index.Levels[i].Refresh()
	}
	if err := s.values.Drop(pos); err != nil {
		return Series{}, fmt.Errorf("Series.Drop(): %v", err)
	}
	return s, nil
}

// Append adds a row at the end and returns a new Series.
func (s Series) Append(val interface{}, idx []interface{}) Series {
	s, _ = s.Insert(s.Len(), val, idx)
	return s
}
