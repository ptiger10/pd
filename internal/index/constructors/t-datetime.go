package constructors

import (
	"time"

	"github.com/ptiger10/pd/internal/index"
	constructVal "github.com/ptiger10/pd/internal/values/constructors"
	"github.com/ptiger10/pd/kinds"
)

// SliceDateTime converts []time.Time{} -> IndexLevel of kind reflect.Struct
func SliceDateTime(data []time.Time, name string) index.Level {
	level := level(
		constructVal.SliceDateTime(data),
		kinds.DateTime,
		name,
	)
	return level
}
