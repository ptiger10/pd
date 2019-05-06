package values

import "github.com/ptiger10/pd/kinds"

// The Values interface is the primary means of handling a collection of values
// Thes same interface and value types are used for both Series values and Index labels
type Values interface {
	Len() int                  // number of value/null structs
	All() []interface{}        // all value elements regardless of null status, as interface slice
	Vals() interface{}         // a slice of values elements in their native form, ready for type assertion
	In([]int) Values           // Value/Null elements at one or more integer positions
	Valid() []int              // integer positions of non-null values
	Null() []int               // integer positions of null values
	Element(int) []interface{} // value element at an integer position, in the form []interface{} {val, null}

	ToFloat() Values
	ToInt() Values
	ToString() Values
	ToBool() Values
	ToDateTime() Values
	ToInterface() Values
}

// Factory is an extended representation of values containing values and kind
type Factory struct {
	V    Values
	Kind kinds.Kind
}
