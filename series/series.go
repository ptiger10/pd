package series

import (
	"reflect"
)

// Series Type options
const (
	Float    = reflect.Float64
	Int      = reflect.Int64
	String   = reflect.String
	Bool     = reflect.Bool
	DateTime = reflect.Struct        // time.Time{} are the only structs accepted by constructor
	None     = reflect.UnsafePointer // pseudo-nil value for type reflect.Kind
)

// A Series is a 1-D data container with a labeled index, static type, and the ability to handle null values
type Series struct {
	Index  Index
	Values Values
	Kind   reflect.Kind
	Name   string
}

type Values interface {
	describe() string
	count() int
}
