// Package series defines the Series, a typed 1-dimensional data structure with an n-level index, analogous to a column in a spreadsheet.
package series

import (
	"fmt"
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
	Index    Index
	InPlace  InPlace
}

func (s *Series) String() string {
	if Equal(s, newEmptySeries()) {
		return "{Empty Series}"
	}
	return s.print()
}

// InPlace contains methods for modifying a Series in place.
type InPlace struct {
	s *Series
}

func (ip InPlace) String() string {
	printer := "{InPlace Series Handler}\n"
	printer += "Methods:\n"
	t := reflect.TypeOf(InPlace{})
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		printer += fmt.Sprintln(method.Name)
	}
	return printer
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
		{"Value", el.Value},
		{"Null", el.Null},
		{"Labels", el.Labels},
		{"LabelTypes", el.LabelTypes},
	} {
		// LabelTypes is 10 characters wide, so left padding set to 10
		printStr += fmt.Sprintf("%10v:%v%v\n", pair[0], strings.Repeat(" ", values.GetDisplayElementWhitespaceBuffer()), pair[1])
	}
	return printStr
}

// The Config struct can be used in the custom Series constructor to name the Series or specify its data type.
type Config struct {
	Name            string
	DataType        options.DataType
	Index           interface{}
	IndexName       string
	MultiIndex      []interface{}
	MultiIndexNames []string
	Manual          bool
}

// A Grouping returns a collection of index labels with mutually exclusive integer positions.
type Grouping struct {
	s      *Series
	groups map[string]*group
}

func (g Grouping) String() string {
	printer := fmt.Sprintf("{Series Grouping | NumGroups: %v, Groups: [%v]}\n", len(g.groups), strings.Join(g.Groups(), ", "))
	return printer
}

// Index contains index selection and conversion
type Index struct {
	s *Series
}

func (idx Index) String() string {
	printer := fmt.Sprintf("{Series Index | Len: %d, NumLevels: %d}\n", idx.Len(), idx.s.NumLevels())
	return printer
}
