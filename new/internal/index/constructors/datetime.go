package constructors

import (
	"github.com/ptiger10/pd/new/internal/index"
	constructVal "github.com/ptiger10/pd/new/internal/values/constructors"
	"github.com/ptiger10/pd/new/kinds"
)

// SliceDateTime converts []time.Time{} -> IndexLevel of kind reflect.Struct
func SliceDateTime(data interface{}, name string) index.Level {
	level := level(
		constructVal.SliceDateTime(data),
		kinds.DateTime,
		name,
	)
	return level
}
