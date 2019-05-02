package values

import (
	"fmt"
	"math"

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

// [START Converters]

// ToFloat converts IntValues to FloatValues
//
// 1: 1.0
func (vals IntValues) ToFloat() Values {
	var ret FloatValues
	for _, val := range vals {
		if val.Null {
			ret = append(ret, Float(math.NaN(), true))
		} else {
			v := float64(val.V)
			ret = append(ret, Float(v, false))
		}
	}
	return ret
}

// ToInt returns itself
func (vals IntValues) ToInt() Values {
	return vals
}

// [END Converters]

// [START Methods]

// Describe the values in the collection
func (vals IntValues) Describe() string {
	offset := options.DisplayValuesWhitespaceBuffer
	l := len(vals)
	len := fmt.Sprintf("%-*s%d\n", offset, "len", l)
	return fmt.Sprint(len)
}

// [END Methods]
