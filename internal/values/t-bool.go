package values

import (
	"math"
	"time"
)

// [START Definitions]

// BoolValues is a slice of bool-typed Value/Null structs
type BoolValues []BoolValue

// A BoolValue is one bool-typed Value/Null struct
type BoolValue struct {
	v    bool
	null bool
}

// Bool constructs a BoolValue
func Bool(v bool, null bool) BoolValue {
	return BoolValue{
		v:    v,
		null: null,
	}
}

// [END Definitions]

// [START Converters]

// ToFloat converts BoolValues to FloatValues
//
// true: 1.0, false: 0.0, null: NaN
func (vals BoolValues) ToFloat() Values {
	var ret FloatValues
	for _, val := range vals {
		if val.null {
			ret = append(ret, Float(math.NaN(), true))
		} else if val.v {
			ret = append(ret, Float(1, false))
		} else {
			ret = append(ret, Float(0, false))
		}
	}
	return ret
}

// ToInt converts BoolValues to IntValues
//
// true: 1, false: 0, null: 0
func (vals BoolValues) ToInt() Values {
	var ret IntValues
	for _, val := range vals {
		if val.null {
			ret = append(ret, Int(0, true))
		} else if val.v {
			ret = append(ret, Int(1, false))
		} else {
			ret = append(ret, Int(0, false))
		}
	}
	return ret
}

// ToBool returns itself
func (vals BoolValues) ToBool() Values {
	return vals
}

// ToDateTime converts BoolValues to DateTimeValues
//
// notnull: time.Date(1970,1,1,0,0,0,0,time.UTC)
func (vals BoolValues) ToDateTime() Values {
	epochDate := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	var ret DateTimeValues
	for _, val := range vals {
		if val.null {
			ret = append(ret, DateTime(time.Time{}, true))
		} else {
			ret = append(ret, DateTime(epochDate, false))
		}
	}
	return ret
}

// [END Converters]

// [START Methods]

// Describe the values in the collection

// [END Methods]
