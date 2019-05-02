package constructors

import (
	"github.com/ptiger10/pd/new/internal/index"
	constructVal "github.com/ptiger10/pd/new/internal/values/constructors"
	"github.com/ptiger10/pd/new/kinds"
)

// SliceString converts []string -> IndexLevel of kind reflect.String
func SliceString(data interface{}, name string) index.Level {
	level := Level(
		constructVal.SliceString(data),
		kinds.String,
		name,
	)
	return level
}
