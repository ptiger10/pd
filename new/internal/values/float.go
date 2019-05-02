package values

import (
	"fmt"
	"time"

	"github.com/ptiger10/pd/new/options"
)

// [START Definitions]

// FloatValues is a slice of Float64-typed Value/Null structs
type FloatValues []FloatValue

// A FloatValue is one Float64-typed Value/Null struct
type FloatValue struct {
	v    float64
	null bool
}

// Float constructs a FloatValue
func Float(v float64, null bool) FloatValue {
	return FloatValue{
		v:    v,
		null: null,
	}
}

// [END Definitions]

// [START Converters]

// ToFloat returns itself
func (vals FloatValues) ToFloat() Values {
	return vals
}

// ToInt converts FloatValues to IntValues
//
// 1.9: 1, 1.5: 1, null: 0
func (vals FloatValues) ToInt() Values {
	var ret IntValues
	for _, val := range vals {
		if val.null {
			ret = append(ret, Int(0, true))
		} else {
			v := int64(val.v)
			ret = append(ret, Int(v, false))
		}
	}
	return ret
}

// ToBool converts FloatValues to BoolValues
//
// x != 0: true; x == 0: false; null: false
func (vals FloatValues) ToBool() Values {
	var ret BoolValues
	for _, val := range vals {
		if val.null {
			ret = append(ret, Bool(false, true))
		} else {
			if val.v == 0 {
				ret = append(ret, Bool(false, false))
			} else {
				ret = append(ret, Bool(true, false))
			}
		}
	}
	return ret
}

// ToDateTime converts FloatValues to DateTimeValues.
// Tries to convert from Unix EPOCH time, otherwise returns null
func (vals FloatValues) ToDateTime() Values {
	var ret DateTimeValues
	for _, val := range vals {
		if val.null {
			ret = append(ret, DateTime(time.Time{}, true))
		} else {
			ret = append(ret, floatToDateTime(val.v))
		}
	}
	return ret
}

func floatToDateTime(f float64) DateTimeValue {
	v := time.Unix(int64(f), 0)
	if v == (time.Time{}) {
		return DateTime(time.Time{}, true)
	}
	return DateTime(v, false)
}

// [END Converters]

// [START Methods]

// Describe the values in the collection
func (vals FloatValues) Describe() string {
	offset := options.DisplayValuesWhitespaceBuffer
	l := len(vals)
	len := fmt.Sprintf("%-*s%d\n", offset, "len", l)
	return fmt.Sprint(len)
}

// [END Methods]
