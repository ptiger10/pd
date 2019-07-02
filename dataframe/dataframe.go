package dataframe

import (
	"github.com/ptiger10/pd/internal/index"
	"github.com/ptiger10/pd/internal/values"
	"github.com/ptiger10/pd/options"
)

// A DataFrame is a 2D collection of one or more Series with a shared index and associated columns.
type DataFrame struct {
	name    string
	vals    []values.Container
	cols    index.Columns
	index   index.Index
	InPlace InPlace
	Index   Index
}

// A Row is a single row in a DataFrame.
type Row struct {
	Values     []interface{}
	Nulls      []bool
	ValueTypes []options.DataType
	Labels     []interface{}
	LabelTypes []options.DataType
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
}

// A Grouping returns a collection of index labels with mutually exclusive integer positions.
type Grouping struct {
	df     *DataFrame
	groups map[string]*group
}
