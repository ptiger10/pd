package constructors

import (
	"github.com/ptiger10/pd/new/internal/index"
	constructVal "github.com/ptiger10/pd/new/internal/values/constructors"
	"github.com/ptiger10/pd/new/kinds"
)

// SliceInt converts []int64 -> IndexLevel of kind reflect.Int64
func SliceInt(data []int64, name string) index.Level {
	level := level(
		constructVal.SliceInt(data),
		kinds.Int,
		name,
	)
	return level
}
