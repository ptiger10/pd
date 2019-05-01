package values

import (
	"fmt"

	"github.com/ptiger10/pd/new/options"
)

// [START Definitions]

// BoolValues is a slice of bool-typed Value/Null structs
type BoolValues []BoolValue

// A BoolValue is one bool-typed Value/Null struct
type BoolValue struct {
	V    bool
	Null bool
}

// Bool constructs a BoolValue
func Bool(v bool, null bool) BoolValue {
	return BoolValue{
		V:    v,
		Null: null,
	}
}

// [END Definitions]

// [START Methods]

// Describe the values in the collection
func (vals BoolValues) Describe() string {
	offset := options.DisplayValuesWhitespaceBuffer
	l := len(vals)
	len := fmt.Sprintf("%-*s%d\n", offset, "len", l)
	return fmt.Sprint(len)
}

// [END Methods]
