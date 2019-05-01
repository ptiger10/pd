package values

import (
	"fmt"

	"github.com/ptiger10/pd/new/options"
)

// [START Definitions]

// StringValues is a slice of string-typed Value/Null structs
type StringValues []StringValue

// A StringValue is one string-typed Value/Null struct
type StringValue struct {
	V    string
	Null bool
}

// String constructs a StringValue
func String(v string, null bool) StringValue {
	return StringValue{
		V:    v,
		Null: null,
	}
}

// [END Definitions]

// [START Methods]

// Describe the values in the collection
func (vals StringValues) Describe() string {
	offset := options.DisplayValuesWhitespaceBuffer
	l := len(vals)
	len := fmt.Sprintf("%-*s%d\n", offset, "len", l)
	return fmt.Sprint(len)
}

// [END Methods]
