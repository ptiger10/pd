package constructors

import (
	"github.com/ptiger10/pd/new/internal/index"
	constructVal "github.com/ptiger10/pd/new/internal/values/constructors"
	"github.com/ptiger10/pd/new/kinds"
)

// SliceFloat converts []float (of any variety) -> IndexLevel of kind reflect.Float64
func SliceFloat(data interface{}, name string) index.Level {
	level := Level(
		constructVal.SliceFloat(data),
		kinds.Float,
		name,
	)
	return level
}
