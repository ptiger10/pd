package index

import constructVal "github.com/ptiger10/pd/internal/values/constructors"

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
	defaultRange := makeRange(0, length)
	factory := constructVal.SliceInt(defaultRange)
	level := newLevel(factory.V, factory.Kind, "")
	return New(level)
}

// makeRange returns a sequential series of numbers for use in default constructors
func makeRange(min, max int) []int64 {
	a := make([]int64, max-min)
	for i := range a {
		a[i] = int64(min + i)
	}
	return a
}
