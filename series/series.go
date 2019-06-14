package series

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/ptiger10/pd/internal/index"
	"github.com/ptiger10/pd/internal/values"
	"github.com/ptiger10/pd/options"
)

// A Series is a 1-D data container with a labeled index, static type, and the ability to handle null values
type Series struct {
	index    index.Index
	values   values.Values
	datatype options.DataType
	Name     string
	Apply    Apply
	Filter   Filter
	Index    Index
	InPlace  InPlace
}

// An Element is a single item in a Series.
type Element struct {
	Value      interface{}
	Null       bool
	Labels     []interface{}
	LabelKinds []options.DataType
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
		printStr += fmt.Sprintf("%10v:%v%v\n", pair[0], strings.Repeat(" ", options.GetDisplayElementWhitespaceBuffer()), pair[1])
	}
	return printStr
}

// Element returns information about the value and index labels at this position.
func (s *Series) Element(position int) Element {
	elem := s.values.Element(position)

	var idx []interface{}
	for _, lvl := range s.index.Levels {
		idxElem := lvl.Labels.Element(position)
		idxVal := idxElem.Value
		idx = append(idx, idxVal)
	}
	idxKinds := s.index.Kinds()
	return Element{elem.Value, elem.Null, idx, idxKinds}
}

// DataType is the data type of the Series' values. Mimics reflect.Type with the addition of time.Time as DateTime.
func (s *Series) DataType() string {
	return fmt.Sprint(s.datatype)
}

// Copy creates a new deep copy of a Series.
func (s *Series) Copy() *Series {
	idx := s.index.Copy()
	valsCopy := s.values.Copy()
	copyS := &Series{
		values:   valsCopy,
		index:    idx,
		datatype: s.datatype,
		Name:     s.Name,
	}
	copyS.Apply = Apply{s: copyS}
	copyS.Filter = Filter{s: copyS}
	copyS.Index = Index{s: copyS}
	copyS.InPlace = InPlace{s: copyS}
	return copyS
}

// in copies a Series then subsets it to include only index items and values at the positions supplied
func (s *Series) in(positions []int) (*Series, error) {
	if err := s.ensureAlignment(); err != nil {
		return s, fmt.Errorf("Series.in(): %v", err)
	}
	if positions == nil {
		return nil, nil
	}

	newS := s.Copy()
	values, err := newS.values.In(positions)
	if err != nil {
		return nil, fmt.Errorf("Series.in() values: %v", err)
	}
	newS.values = values
	for i, level := range newS.index.Levels {
		newS.index.Levels[i].Labels, err = level.Labels.In(positions)
		if err != nil {
			return nil, fmt.Errorf("Series.in() index: %v", err)
		}
	}
	newS.index.Refresh()
	return newS, nil
}

func (s *Series) mustIn(positions []int) *Series {
	s, err := s.in(positions)
	if err != nil {
		log.Printf("Internal error: %v\n", err)
		return nil
	}
	return s
}

func seriesEquals(s1, s2 *Series) bool {
	sameIndex := reflect.DeepEqual(s1.index, s2.index)
	sameValues := reflect.DeepEqual(s1.values, s2.values)
	sameName := s1.Name == s2.Name
	sameKind := s1.datatype == s2.datatype
	if sameIndex && sameValues && sameName && sameKind {
		return true
	}
	return false
}

// Len returns the number of Elements (i.e., Value/Null pairs) in the Series.
func (s *Series) Len() int {
	if s.values == nil {
		return 0
	}
	return s.values.Len()
}

// valid returns integer positions of valid (i.e., non-null) values in the series.
func (s *Series) valid() []int {
	var ret []int
	for i := 0; i < s.Len(); i++ {
		if !s.values.Element(i).Null {
			ret = append(ret, i)
		}
	}
	return ret
}

// null returns the integer position of all null values in the collection.
func (s *Series) null() []int {
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
func (s *Series) all() []interface{} {
	var ret []interface{}
	for i := 0; i < s.Len(); i++ {
		ret = append(ret, s.values.Element(i).Value)
	}
	return ret
}

func (s *Series) replace(s2 *Series) {
	s.Name = s2.Name
	s.datatype = s2.datatype
	s.values = s2.values
	s.index = s2.index
}
