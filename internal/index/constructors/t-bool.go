package constructors

import (
	"github.com/ptiger10/pd/internal/index"
	constructVal "github.com/ptiger10/pd/internal/values/constructors"
	"github.com/ptiger10/pd/kinds"
)

// SliceBool converts []bool -> IndexLevel of kind reflect.Bool
func SliceBool(data []bool, name string) index.Level {
	level := level(
		constructVal.SliceBool(data),
		kinds.Bool,
		name,
	)
	return level
}
