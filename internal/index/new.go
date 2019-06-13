package index

import (
	"github.com/ptiger10/pd/options"
	"github.com/ptiger10/pd/internal/values"
)

// New receives one or more Levels and returns a new Index.
// Expects that Levels already have .LabelMap and .Longest set.
func New(levels ...Level) Index {
	idx := Index{
		Levels: levels,
	}
	idx.UpdateNameMap()
	return idx
}

// Default creates an index with one unnamed index level and range labels (0, 1, 2, ...n)
func Default(length int) Index {
	level := DefaultLevel(length, "")
	return New(level)
}

// DefaultLevel creates an unnamed index level with range labels (0, 1, 2, ...n)
func DefaultLevel(n int, name string) Level {
	v := values.NewDefault(n)
	level := newLevel(v, options.Int64, name)
	return level
}
