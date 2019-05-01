package values

import (
	"fmt"

	"github.com/ptiger10/pd/new/options"
)

// [START Definitions]

// IntValues is a slice of Int64-typed Value/Null structs
type IntValues []IntValue

// An IntValue is one Int64-typed Value/Null struct
type IntValue struct {
	V    int64
	Null bool
}

// Int constructs an IntValue
func Int(v int64, null bool) IntValue {
	return IntValue{
		V:    v,
		Null: null,
	}
}

// [END Definitions]

// [START Methods]

// Describe the values in the collection
func (vals IntValues) Describe() string {
	offset := options.DisplayValuesWhitespaceBuffer
	l := len(vals)
	len := fmt.Sprintf("%-*s%d\n", offset, "len", l)
	return fmt.Sprint(len)
}

// [END Methods]
