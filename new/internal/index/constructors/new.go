package constructors

import "github.com/ptiger10/pd/new/internal/index"

// New receives one or more Levels and returns a new Index
func New(levels []index.Level) index.Index {
	return index.Index{Levels: levels}
}
