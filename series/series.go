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

// An Element is a single item in a Series.
type Element struct {
	Value      interface{}
	Null       bool
	Labels     []interface{}
	LabelTypes []options.DataType
}
