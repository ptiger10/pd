package values

import (
	"fmt"
	"math"
	"strconv"

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

// [START Converters]

// ToFloat converts StringValues to FloatValues
//
// "1": 1.0, Null: NaN
func (vals StringValues) ToFloat() Values {
	var ret FloatValues
	for _, val := range vals {
		if val.Null {
			ret = append(ret, Float(math.NaN(), true))
		} else {
			v, err := strconv.ParseFloat(val.V, 64)
			if err != nil {
				ret = append(ret, Float(math.NaN(), true))
			} else {
				ret = append(ret, Float(v, false))
			}
		}
	}
	return ret
}

// ToInt converts StringValues to IntValues
//
// "1": 1, null: NaN
func (vals StringValues) ToInt() Values {
	var ret IntValues
	for _, val := range vals {
		if val.Null {
			ret = append(ret, Int(0, true))
		} else {
			v, err := strconv.Atoi(val.V)
			if err != nil {
				ret = append(ret, Int(0, true))
			} else {
				ret = append(ret, Int(int64(v), false))
			}
		}
	}
	return ret
}

// ToBool converts StringValues to BoolValues
//
// null: false; notnull: true
func (vals StringValues) ToBool() Values {
	var ret BoolValues
	for _, val := range vals {
		if val.Null {
			ret = append(ret, Bool(false, true))
		} else {
			ret = append(ret, Bool(true, false))
		}
	}
	return ret
}

// [END Converters]

// [START Methods]

// Describe the values in the collection
func (vals StringValues) Describe() string {
	offset := options.DisplayValuesWhitespaceBuffer
	l := len(vals)
	len := fmt.Sprintf("%-*s%d\n", offset, "len", l)
	return fmt.Sprint(len)
}

// [END Methods]
