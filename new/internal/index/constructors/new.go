package constructors

import "github.com/ptiger10/pd/new/internal/index"

// New receives one or more Levels and returns a new Index
func New(levels []index.Level) index.Index {
	idx := index.Index{
		Levels:  levels,
		NameMap: make(index.LabelMap),
	}
	for i, level := range levels {
		idx.NameMap[level.Name] = append(idx.NameMap[level.Name], i)
		idx.Levels[i].ComputeLongest()
	}
	return idx
}

// Default creates an unnamed index with range labels (0, 1, 2, ...n)
func Default(length int) index.Index {
	defaultRange := makeRange(0, length)
	level := SliceInt(defaultRange)
	level.ComputeLongest()
	idx := New([]index.Level{level})
	return idx
}

// makeRange returns a sequential series of numbers for use in default constructors
func makeRange(min, max int) []int64 {
	a := make([]int64, max-min)
	for i := range a {
		a[i] = int64(min + i)
	}
	return a
}
