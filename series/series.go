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
	name     string
	Apply    Apply
	Index    Index
	InPlace  InPlace
}

// An Element is a single item in a Series.
type Element struct {
	Value      interface{}
	Null       bool
	Labels     []interface{}
	LabelTypes []options.DataType
}

func (el Element) String() string {
	var printStr string
	for _, pair := range [][]interface{}{
		[]interface{}{"Value", el.Value},
		[]interface{}{"Null", el.Null},
		[]interface{}{"Labels", el.Labels},
		[]interface{}{"LabelTypes", el.LabelTypes},
	} {
		// LabelTypes is 10 characters wide, so left padding set to 10
		printStr += fmt.Sprintf("%10v:%v%v\n", pair[0], strings.Repeat(" ", options.GetDisplayElementWhitespaceBuffer()), pair[1])
	}
	return printStr
}

// Element returns information about the value and index labels at this position.
func (s *Series) Element(position int) Element {
	elem := s.values.Element(position)
	idxElems := s.index.Elements(position)
	return Element{elem.Value, elem.Null, idxElems.Labels, idxElems.DataTypes}
}

// DataType is the data type of the Series' values. Mimics reflect.Type with the addition of time.Time as DateTime.
func (s *Series) DataType() string {
	return fmt.Sprint(s.datatype)
}

// selectByRows copies a Series then subsets it to include only index items and values at the positions supplied
func (s *Series) selectByRows(positions []int) (*Series, error) {
	if err := s.ensureAlignment(); err != nil {
		return s, fmt.Errorf("series internal alignment error: %v", err)
	}
	if positions == nil {
		return newEmptySeries(), nil
	}

	s = s.Copy()
	values, err := s.values.In(positions)
	if err != nil {
		return newEmptySeries(), fmt.Errorf("series.selectByRows() selecting rows: %v", err)
	}
	s.values = values
	idx, err := s.index.In(positions)
	if err != nil {
		return newEmptySeries(), fmt.Errorf("series.selectByRows(): %v", err)
	}
	s.index = idx
	return s, nil
}

func (s *Series) mustSelectRows(positions []int) *Series {
	s, err := s.selectByRows(positions)
	if err != nil {
		log.Printf("Internal error: %v\n", err)
		return newEmptySeries()
	}
	return s
}

// Equal compares whether two series are equivalent.
func Equal(s1, s2 *Series) bool {
	if !reflect.DeepEqual(s1.values, s2.values) {
		return false
	}
	if !reflect.DeepEqual(s1.index, s2.index) {
		return false
	}
	if s1.name != s2.name {
		return false
	}
	if s1.datatype != s2.datatype {
		return false
	}
	return true
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
	s.name = s2.name
	s.datatype = s2.datatype
	s.values = s2.values
	s.index = s2.index
}

// MaxWidth returns the max number of characters in any value in the Series.
// For use in printing Series and DataFrames.
func (s *Series) MaxWidth() int {
	var max int
	for _, v := range s.all() {
		if length := len(fmt.Sprint(v)); length > max {
			max = length
		}
	}
	return max
}

// Name returns the Series' name.
func (s Series) Name() string {
	return s.name
}
