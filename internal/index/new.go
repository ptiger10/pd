package index

import "github.com/ptiger10/pd/internal/values"

// New receives one or more Levels and returns a new Index.
// Expects that Levels already have .LabelMap and .Longest set.
func New(levels ...Level) Index {
	idx := Index{
		Levels: levels,
	}
	idx.UpdateNameMap()
	return idx
}

// Default creates an unnamed index level with range labels (0, 1, 2, ...n)
func Default(length int) Index {
	defaultRange := values.MakeRange(0, length)
	factory := values.NewSliceInt(defaultRange)
	level := newLevel(factory.Values, factory.DataType, "")
	return New(level)
}
