package values

import (
	"fmt"

	"github.com/ptiger10/pd/new/options"
)

// [START Definitions]

// InterfaceValues is a slice of interface-typed Value/Null structs
type InterfaceValues []InterfaceValue

// An InterfaceValue is one interface-typed Value/Null struct
type InterfaceValue struct {
	V    interface{}
	Null bool
}

// Interface constructs an InterfaceValue
func Interface(v interface{}, null bool) InterfaceValue {
	return InterfaceValue{
		V:    v,
		Null: null,
	}
}

// [END Definitions]

// [START Methods]

// Describe the values in the collection
func (vals InterfaceValues) Describe() string {
	offset := options.DisplayValuesWhitespaceBuffer
	l := len(vals)
	len := fmt.Sprintf("%-*s%d\n", offset, "len", l)
	return fmt.Sprint(len)
}

// [END Methods]
