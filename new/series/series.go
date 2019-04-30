package series

import (
	"reflect"

	"github.com/ptiger10/pd/new/internal/index"
	"github.com/ptiger10/pd/new/internal/values"
)

// A Series is a 1-D data container with a labeled index, static type, and the ability to handle null values
type Series struct {
	Index  index.Index
	Values values.Values
	Kind   reflect.Kind
	Name   string
}
