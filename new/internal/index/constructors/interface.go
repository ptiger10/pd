package constructors

import (
	"github.com/ptiger10/pd/new/internal/index"
	constructVal "github.com/ptiger10/pd/new/internal/values/constructors"
	"github.com/ptiger10/pd/new/kinds"
)

// SliceInterface converts []interface{} -> IndexLevel of kind reflect.Interface
func SliceInterface(data interface{}, name string) index.Level {
	level := level(
		constructVal.SliceInterface(data),
		kinds.Interface,
		name,
	)
	return level
}
