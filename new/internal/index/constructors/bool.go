package constructors

import (
	"github.com/ptiger10/pd/new/internal/index"
	constructVal "github.com/ptiger10/pd/new/internal/values/constructors"
	"github.com/ptiger10/pd/new/kinds"
)

// SliceBool converts []bool -> IndexLevel of kind reflect.Bool
func SliceBool(data interface{}, name string) index.Level {
	level := level(
		constructVal.SliceBool(data),
		kinds.Bool,
		name,
	)
	return level
}
