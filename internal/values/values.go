package values

import (
	"github.com/ptiger10/pd/kinds"
)

// The Values interface is the primary means of handling a collection of values
// Thes same interface and value types are used for both Series values and Index labels
type Values interface {
	Len() int                      // number of Value/Null structs
	Vals() interface{}             // a slice of values in their native form, ready for type assertion
	In([]int) (Values, error)      // a new Values object comprised of the Value/Null pairs at one or more integer positions
	Element(int) Elem              // Value/Null pair at an integer position
	Set(int, interface{}) error    // overwrite the value/null struct at an integer position
	Copy() Values                  // clone the Values
	Insert(int, interface{}) error // insert a Value/Null pair at an integer position
	Drop(int) error                // drop a Value/Null pair at an integer position

	ToFloat() Values
	ToInt() Values
	ToString() Values
	ToBool() Values
	ToDateTime() Values
	ToInterface() Values
}

// Factory contains Values (a list of Value/Null pairs satisfying the Values interface) and Kind.
type Factory struct {
	Values Values
	Kind   kinds.Kind
}

// An Elem is a single Value/Null pair.
type Elem struct {
	Value interface{}
	Null  bool
}
