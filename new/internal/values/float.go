package values

import (
	"fmt"

	"github.com/ptiger10/pd/new/options"
)

// [START Definitions]

// FloatValues is a slice of Float64-typed Value/Null structs
type FloatValues []FloatValue

// A FloatValue is one Float64-typed Value/Null struct
type FloatValue struct {
	V    float64
	Null bool
}

// Float constructs a FloatValue
func Float(v float64, null bool) FloatValue {
	return FloatValue{
		V:    v,
		Null: null,
	}
}

// [END Definitions]

// [START Methods]

// Describe the values in the collection
func (vals FloatValues) Describe() string {
	offset := options.DisplayValuesWhitespaceBuffer
	l := len(vals)
	len := fmt.Sprintf("%-*s%d\n", offset, "len", l)
	return fmt.Sprint(len)
}

// [END Methods]
