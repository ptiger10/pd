package values

import (
	"fmt"
	"math"
	"time"

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

// ToBool converts IntValues to BoolValues
//
// x != 0: true; x == 0: false; null: false
func (vals IntValues) ToBool() Values {
	var ret BoolValues
	for _, val := range vals {
		if val.Null {
			ret = append(ret, Bool(false, true))
		} else {
			if val.V == 0 {
				ret = append(ret, Bool(false, false))
			} else {
				ret = append(ret, Bool(true, false))
			}
		}
	}
	return ret
}

// ToDateTime converts IntValues to DateTimeValues.
// Tries to convert from Unix EPOCH time, otherwise returns null
func (vals IntValues) ToDateTime() Values {
	var ret DateTimeValues
	for _, val := range vals {
		if val.Null {
			ret = append(ret, DateTime(time.Time{}, true))
		} else {
			ret = append(ret, intToDateTime(val.V))
		}
	}
	return ret
}

func intToDateTime(i int64) DateTimeValue {
	v := time.Unix(i, 0)
	if v == (time.Time{}) {
		return DateTime(time.Time{}, true)
	}
	return DateTime(v, false)
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
