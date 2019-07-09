package dataframe

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/ptiger10/pd/internal/index"
	"github.com/ptiger10/pd/internal/values"
	"github.com/ptiger10/pd/options"
)

// A DataFrame is a 2D collection of one or more Series with a shared index and associated columns.
type DataFrame struct {
	name    string
	vals    []values.Container
	cols    index.Columns
	Columns Columns
	index   index.Index
	Index   Index
	InPlace InPlace
}

func (df *DataFrame) String() string {
	if Equal(df, newEmptyDataFrame()) {
		return "{Empty DataFrame}"
	}
	return df.print()
}

// Index contains index level data.
type Index struct {
	df *DataFrame
}

func (idx Index) String() string {
	printer := fmt.Sprintf("{DataFrame Index | Len: %d, NumLevels: %d}\n", idx.Len(), idx.df.IndexLevels())
	return printer
}

// Columns contains column level data.
type Columns struct {
	df *DataFrame
}

func (col Columns) String() string {
	printer := fmt.Sprintf("{DataFrame Columns | NumCols: %d, NumLevels: %d}\n", col.df.NumCols(), col.df.ColLevels())
	return printer
}

// A Row is a single row in a DataFrame.
type Row struct {
	Values     []interface{}
	Nulls      []bool
	ValueTypes []options.DataType
	Labels     []interface{}
	LabelTypes []options.DataType
}

func (r Row) String() string {
	var printStr string
	for _, pair := range [][]interface{}{
		[]interface{}{"Values", r.Values},
		[]interface{}{"IsNull", r.Nulls},
		[]interface{}{"ValueTypes", r.ValueTypes},
		[]interface{}{"Labels", r.Labels},
		[]interface{}{"LabelTypes", r.LabelTypes},
	} {
		// LabelTypes is 10 characters wide, so left padding set to 10
		printStr += fmt.Sprintf("%10v:%v%v\n", pair[0], strings.Repeat(" ", values.GetDisplayElementWhitespaceBuffer()), pair[1])
	}
	return printStr
}

// Config customizes the DataFrame constructor.
type Config struct {
	Name            string
	DataType        options.DataType
	Index           interface{}
	IndexName       string
	MultiIndex      []interface{}
	MultiIndexNames []string
	Col             []string
	ColName         string
	MultiCol        [][]string
	MultiColNames   []string
	Manual          bool
}

// A Grouping returns a collection of index labels with mutually exclusive integer positions.
type Grouping struct {
	df     *DataFrame
	groups map[string]*group
	err    bool
}

func (g Grouping) String() string {
	printer := fmt.Sprintf("{DataFrame Grouping | NumGroups: %v, Groups: [%v]}\n", len(g.groups), strings.Join(g.Groups(), ", "))
	return printer
}

// InPlace contains methods for modifying a DataFrame in place.
type InPlace struct {
	df *DataFrame
}

func (ip InPlace) String() string {
	printer := "{InPlace DataFrame Handler}\n"
	printer += "Methods:\n"
	t := reflect.TypeOf(InPlace{})
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		printer += fmt.Sprintln(method.Name)
	}
	return printer
}
