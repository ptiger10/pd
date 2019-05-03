package constructors

import (
	"github.com/ptiger10/pd/internal/index"
	constructVal "github.com/ptiger10/pd/internal/values/constructors"
	"github.com/ptiger10/pd/kinds"
)

// SliceString converts []string -> IndexLevel of kind reflect.String
func SliceString(data []string, name string) index.Level {
	level := level(
		constructVal.SliceString(data),
		kinds.String,
		name,
	)
	return level
}
