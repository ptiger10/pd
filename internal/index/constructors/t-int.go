package constructors

import (
	"github.com/ptiger10/pd/internal/index"
	constructVal "github.com/ptiger10/pd/internal/values/constructors"
	"github.com/ptiger10/pd/kinds"
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
