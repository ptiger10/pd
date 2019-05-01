package constructors

import (
	"reflect"

	"github.com/ptiger10/pd/new/internal/index"
	"github.com/ptiger10/pd/new/internal/values"
)

// Level returns an Index level with updated label map and longest value computed.
// NB: Create labels using the values.constructors methods
func Level(labels values.Values, kind reflect.Kind, name string) index.Level {
	lvl := index.Level{
		Labels: labels,
		Kind:   kind,
		Name:   name,
	}
	lvl.UpdateLabelMap()
	lvl.ComputeLongest()
	return lvl
}
