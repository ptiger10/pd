package constructors

import (
	"github.com/ptiger10/pd/internal/index"
	constructVal "github.com/ptiger10/pd/internal/values/constructors"
	"github.com/ptiger10/pd/kinds"
)

// SliceInterface converts []interface{} -> IndexLevel of kind reflect.Interface
func SliceInterface(data []interface{}, name string) index.Level {
	level := level(
		constructVal.SliceInterface(data),
		kinds.Interface,
		name,
	)
	return level
}
