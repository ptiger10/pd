package constructors

import (
	"github.com/ptiger10/pd/internal/index"
	constructVal "github.com/ptiger10/pd/internal/values/constructors"
	"github.com/ptiger10/pd/kinds"
)

// SliceFloat converts []float64 -> IndexLevel of kind reflect.Float64
func SliceFloat(data []float64, name string) index.Level {
	level := level(
		constructVal.SliceFloat(data),
		kinds.Float,
		name,
	)
	return level
}
