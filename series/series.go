package series

import (
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

// InPlace contains methods for modifying a Series in place.
type InPlace struct {
	s *Series
}

// An Element is a single item in a Series.
type Element struct {
	Value      interface{}
	Null       bool
	Labels     []interface{}
	LabelTypes []options.DataType
}

// The Config struct can be used in the custom Series constructor to name the Series or specify its data type.
// Basic usage: New("foo", series.Config{Name: "bar"})
type Config struct {
	Name            string
	DataType        options.DataType
	Index           interface{}
	IndexName       string
	MultiIndex      []interface{}
	MultiIndexNames []string
}

// A Grouping returns a collection of index labels with mutually exclusive integer positions.
type Grouping struct {
	s      *Series
	groups map[string]*group
}
