package series

import (
	"reflect"

	"github.com/ptiger10/pd/new/internal/index"
)

// A Series is a 1-D data container with a labeled index, static type, and the ability to handle null values
type Series struct {
	Index  index.Index
	Values Values
	Kind   reflect.Kind
	Name   string
}

type Values interface {
	Describe() string
	In([]int) interface{}
}
