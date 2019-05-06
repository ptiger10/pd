package values

import (
	"math"
	"time"

	"github.com/ptiger10/pd/kinds"
)

// [START Constructor Functions]

// SliceInt converts []int (of any variety) -> Factory with int64Values
func SliceInt(vals []int64) Factory {
	var v int64Values
	for _, val := range vals {
		v = append(v, int64Val(val, false))
	}
	return Factory{v, kinds.Int}
}

// [END Constructor Functions]

// [START Converters]

// ToFloat converts int64Values to float64Values
//
// 1: 1.0
func (vals int64Values) ToFloat() Values {
	var ret float64Values
	for _, val := range vals {
		if val.null {
			ret = append(ret, float64Val(math.NaN(), true))
		} else {
			v := float64(val.v)
			ret = append(ret, float64Val(v, false))
		}
	}
	return ret
}

// ToInt returns itself
func (vals int64Values) ToInt() Values {
	return vals
}

// ToBool converts int64Values to boolValues
//
// x != 0: true; x == 0: false; null: false
func (vals int64Values) ToBool() Values {
	var ret boolValues
	for _, val := range vals {
		if val.null {
			ret = append(ret, boolVal(false, true))
		} else {
			if val.v == 0 {
				ret = append(ret, boolVal(false, false))
			} else {
				ret = append(ret, boolVal(true, false))
			}
		}
	}
	return ret
}

// ToDateTime converts int64Values to dateTimeValues.
// Tries to convert from Unix EPOCH timestamp.
// Defaults to 1970-01-01 00:00:00 +0000 UTC.
func (vals int64Values) ToDateTime() Values {
	var ret dateTimeValues
	for _, val := range vals {
		if val.null {
			ret = append(ret, dateTimeVal(time.Time{}, true))
		} else {
			ret = append(ret, intToDateTime(val.v))
		}
	}
	return ret
}

func intToDateTime(i int64) dateTimeValue {
	// convert from nanoseconds to seconds
	i /= 1000000000
	v := time.Unix(i, 0).UTC()
	return dateTimeVal(v, false)
}

// [END Converters]
