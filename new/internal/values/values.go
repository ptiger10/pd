package values

import (
	"reflect"
)

type Values interface {
	Describe() string
	// Important - returning []interface{} means you can't type assert values.Values
	In([]int) interface{}
	Kind() reflect.Kind
	// returns all values regardless of null status
	All() []interface{}
}
