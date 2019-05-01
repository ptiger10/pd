package values

import (
	"fmt"
	"time"

	"github.com/ptiger10/pd/new/options"
)

// [START Definitions]

// DateTimeValues is a slice of time.Time-typed Value/Null structs
type DateTimeValues []DateTimeValue

// A DateTimeValue is one time.Time-typed Value/Null struct
type DateTimeValue struct {
	V    time.Time
	Null bool
}

// DateTime constructs a DateTimeValue
func DateTime(v time.Time, null bool) DateTimeValue {
	return DateTimeValue{
		V:    v,
		Null: null,
	}
}

// [END Definitions]

// [START Methods]

// Describe the values in the collection
func (vals DateTimeValues) Describe() string {
	offset := options.DisplayValuesWhitespaceBuffer
	l := len(vals)
	len := fmt.Sprintf("%-*s%d\n", offset, "len", l)
	return fmt.Sprint(len)
}

// [END Methods]
