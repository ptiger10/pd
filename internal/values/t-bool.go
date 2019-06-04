package values

import (
	"math"
	"time"

	"github.com/ptiger10/pd/kinds"
)

// [START Constructor Functions]

// SliceBool converts []bool -> Factory with boolValues
func SliceBool(vals []bool) Factory {
	var v boolValues
	for _, val := range vals {
		v = append(v, boolVal(val, false))
	}
	return Factory{v, kinds.Bool}
}

// [END Constructor Functions]

// [START Converters]
// Set overwrites the

// ToFloat converts boolValues to float64Values.
//
// true: 1.0, false: 0.0, null: NaN
func (vals boolValues) ToFloat() Values {
	var ret float64Values
	for _, val := range vals {
		if val.null {
			ret = append(ret, float64Val(math.NaN(), true))
		} else if val.v {
			ret = append(ret, float64Val(1, false))
		} else {
			ret = append(ret, float64Val(0, false))
		}
	}
	return ret
}

// ToInt converts boolValues to int64Values.
//
// true: 1, false: 0, null: 0
func (vals boolValues) ToInt() Values {
	var ret int64Values
	for _, val := range vals {
		if val.null {
			ret = append(ret, int64Val(0, true))
		} else if val.v {
			ret = append(ret, int64Val(1, false))
		} else {
			ret = append(ret, int64Val(0, false))
		}
	}
	return ret
}

// ToBool returns itself.
func (vals boolValues) ToBool() Values {
	return vals
}

// ToDateTime converts boolValues to dateTimeValues.
//
// notnull: time.Date(1970,1,1,0,0,0,0,time.UTC)
func (vals boolValues) ToDateTime() Values {
	epochDate := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	var ret dateTimeValues
	for _, val := range vals {
		if val.null {
			ret = append(ret, dateTimeVal(time.Time{}, true))
		} else {
			ret = append(ret, dateTimeVal(epochDate, false))
		}
	}
	return ret
}

// [END Converters]
