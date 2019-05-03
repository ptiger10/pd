package values

import (
	"math"
	"time"
)

// [START Definitions]

// DateTimeValues is a slice of time.Time-typed Value/Null structs
type DateTimeValues []DateTimeValue

// A DateTimeValue is one time.Time-typed Value/Null struct
type DateTimeValue struct {
	v    time.Time
	null bool
}

// DateTime constructs a DateTimeValue
func DateTime(v time.Time, null bool) DateTimeValue {
	return DateTimeValue{
		v:    v,
		null: null,
	}
}

// [END Definitions]

// [START Converters]

// ToFloat converts DateTimeValues to FloatValues of the Unix EPOCH timestamp
// (seconds since midnight January 1, 1970)
// 2019-05-01 00:00:00 +0000 UTC: 1556757505
func (vals DateTimeValues) ToFloat() Values {
	var ret FloatValues
	for _, val := range vals {
		if val.null {
			ret = append(ret, Float(math.NaN(), true))
		} else {
			v := val.v.UnixNano()
			ret = append(ret, Float(float64(v), false))
		}
	}
	return ret
}

// ToInt converts DateTimeValues to IntValues of the Unix EPOCH timestamp
// (seconds since midnight January 1, 1970)
//
// 2019-05-01 00:00:00 +0000 UTC: 1556757505
func (vals DateTimeValues) ToInt() Values {
	var ret IntValues
	for _, val := range vals {
		if val.null {
			ret = append(ret, Int(0, true))
		} else {
			v := val.v.UnixNano()
			ret = append(ret, Int(v, false))
		}
	}
	return ret
}

// ToBool converts DateTimeValues to BoolValues
//
// x != time.Time{}: true; x == time.Time{}: false; null: false
func (vals DateTimeValues) ToBool() Values {
	var ret BoolValues
	for _, val := range vals {
		if val.null {
			ret = append(ret, Bool(false, true))
		} else {
			if val.v == (time.Time{}) {
				ret = append(ret, Bool(false, false))
			} else {
				ret = append(ret, Bool(true, false))
			}
		}
	}
	return ret
}

// ToDateTime returns itself
func (vals DateTimeValues) ToDateTime() Values {
	return vals
}

// [END Converters]

// [START Methods]

// Describe the values in the collection
// [END Methods]
