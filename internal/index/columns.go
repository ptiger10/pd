package index

import "github.com/ptiger10/pd/options"

// Columns is a collection of column levels, plus name mappings.
type Columns struct {
	NameMap LabelMap
	Levels  []ColLevel
}

// A ColLevel is a single collection of column labels within a Columns collection, plus label mappings and metadata.
// It is identical to an index Level except for the Labels, which are a simple []interface{} that do not satisfy the values.Values interface.
type ColLevel struct {
	Name     string
	DataType options.DataType
	Labels   []interface{}
	LabelMap LabelMap
}
