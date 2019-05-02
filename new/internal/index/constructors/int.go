package constructors

import (
	"github.com/ptiger10/pd/new/internal/index"
	constructVal "github.com/ptiger10/pd/new/internal/values/constructors"
	"github.com/ptiger10/pd/new/kinds"
)

// SliceInt converts []int (of any variety) -> IndexLevel of kind reflect.Int64
func SliceInt(data interface{}, name string) index.Level {
	level := level(
		constructVal.SliceInt(data),
		kinds.Int,
		name,
	)
	return level
}

// SliceUint converts []uint (of any variety) -> IndexLevel of kind reflect.Int64
func SliceUint(data interface{}, name string) index.Level {
	level := level(
		constructVal.SliceUint(data),
		kinds.Int,
		name,
	)
	return level
}
