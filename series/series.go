package series

import (
	"github.com/ptiger10/pd/internal/index"
	"github.com/ptiger10/pd/internal/values"
	"github.com/ptiger10/pd/kinds"
)

// A Series is a 1-D data container with a labeled index, static type, and the ability to handle null values
type Series struct {
	index  index.Index
	values values.Values
	Kind   kinds.Kind
	Name   string
}
