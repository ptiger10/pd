package series

import (
	"fmt"

	"github.com/ptiger10/pd/internal/index"
	"github.com/ptiger10/pd/internal/values"
	"github.com/ptiger10/pd/kinds"
)

// A Series is a 1-D data container with a labeled index, static type, and the ability to handle null values
type Series struct {
	index  index.Index
	values values.Values
	kind   kinds.Kind
	Name   string
}

// Kind is the data kind of the Series' values. Mimics reflect.Kind with the addition of time.Time
func (s Series) Kind() string {
	return fmt.Sprint(s.kind)
}
