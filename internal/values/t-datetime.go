package values

import (
	"math"
	"time"

	"github.com/ptiger10/pd/kinds"
)

// [START Constructor Functions]

// SliceDateTime converts []time.Time{} -> Factory with dateTimeValues
func SliceDateTime(vals []time.Time) Factory {
	var v dateTimeValues
	for _, val := range vals {
		if (time.Time{}) == val {
			v = append(v, dateTimeVal(val, true))
		} else {
			v = append(v, dateTimeVal(val, false))
		}

	}

	return Factory{v, kinds.DateTime}
}

// [END Constructor Functions]

// [START Converters]

// ToFloat converts dateTimeValues to float64Values of the Unix EPOCH timestamp
// (seconds since midnight January 1, 1970)
// 2019-05-01 00:00:00 +0000 UTC: 1556757505
func (vals dateTimeValues) ToFloat() Values {
	var ret float64Values
	for _, val := range vals {
		if val.null {
			ret = append(ret, float64Val(math.NaN(), true))
		} else {
			v := val.v.UnixNano()
			ret = append(ret, float64Val(float64(v), false))
		}
	}
	return ret
}

// ToInt converts dateTimeValues to int64Values of the Unix EPOCH timestamp
// (seconds since midnight January 1, 1970)
//
// 2019-05-01 00:00:00 +0000 UTC: 1556757505
func (vals dateTimeValues) ToInt() Values {
	var ret int64Values
	for _, val := range vals {
		if val.null {
			ret = append(ret, int64Val(0, true))
		} else {
			v := val.v.UnixNano()
			ret = append(ret, int64Val(v, false))
		}
	}
	return ret
}

// ToBool converts dateTimeValues to boolValues
//
// x != time.Time{}: true; x == time.Time{}: false; null: false
func (vals dateTimeValues) ToBool() Values {
	var ret boolValues
	for _, val := range vals {
		if val.null {
			ret = append(ret, boolVal(false, true))
		} else {
			if val.v == (time.Time{}) {
				ret = append(ret, boolVal(false, false))
			} else {
				ret = append(ret, boolVal(true, false))
			}
		}
	}
	return ret
}

// ToDateTime returns itself
func (vals dateTimeValues) ToDateTime() Values {
	return vals
}

// [END Converters]

// [START Methods]

// Describe the values in the collection
// [END Methods]
