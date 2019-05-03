package constructors

import (
	"github.com/ptiger10/pd/internal/index"
)

// New receives one or more Levels and returns a new Index. Expects that Levels already have labelMappings set.
func New(levels ...index.Level) index.Index {
	idx := index.Index{
		Levels: levels,
	}
	idx.UpdateNameMap()
	return idx
}

// Default creates an unnamed index level with range labels (0, 1, 2, ...n)
func Default(length int) index.Index {
	defaultRange := makeRange(0, length)
	return New(
		SliceInt(defaultRange, ""),
	)
}

// makeRange returns a sequential series of numbers for use in default constructors
func makeRange(min, max int) []int64 {
	a := make([]int64, max-min)
	for i := range a {
		a[i] = int64(min + i)
	}
	return a
}
