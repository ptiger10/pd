package datatypes

import "reflect"

// Type convenience options
const (
	Float    = reflect.Float64
	Int      = reflect.Int64
	String   = reflect.String
	Bool     = reflect.Bool
	DateTime = reflect.Struct        // time.Time{} are the only structs accepted by constructor
	None     = reflect.UnsafePointer // pseudo-nil value for type reflect.Kind
)
