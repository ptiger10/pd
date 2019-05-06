package values

import (
	"math"
	"time"

	"github.com/ptiger10/pd/kinds"
)

// [START Constructor Functions]

// SliceFloat converts []float64  -> Factory with float64Values
func SliceFloat(vals []float64) Factory {
	var v float64Values
	for _, val := range vals {
		if math.IsNaN(val) {
			v = append(v, float64Val(val, true))
			continue
		}
		v = append(v, float64Val(val, false))
	}
	return Factory{v, kinds.Float}
}

// [END Constructor Functions]

// [START Converters]

// ToFloat returns itself
func (vals float64Values) ToFloat() Values {
	return vals
}

// ToInt converts float64Values to int64Values
//
// 1.9: 1, 1.5: 1, null: 0
func (vals float64Values) ToInt() Values {
	var ret int64Values
	for _, val := range vals {
		if val.null {
			ret = append(ret, int64Val(0, true))
		} else {
			v := int64(val.v)
			ret = append(ret, int64Val(v, false))
		}
	}
	return ret
}

// ToBool converts float64Values to boolValues
//
// x != 0: true; x == 0: false; null: false
func (vals float64Values) ToBool() Values {
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

// ToDateTime converts float64Values to dateTimeValues.
// Tries to convert from Unix EPOCH time, otherwise returns null
func (vals float64Values) ToDateTime() Values {
	var ret dateTimeValues
	for _, val := range vals {
		if val.null {
			ret = append(ret, dateTimeVal(time.Time{}, true))
		} else {
			ret = append(ret, floatToDateTime(val.v))
		}
	}
	return ret
}

func floatToDateTime(f float64) dateTimeValue {
	return intToDateTime(int64(f))
}

// [END Converters]
