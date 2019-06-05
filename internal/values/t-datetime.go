package values

import (
	"math"
	"time"

	"github.com/ptiger10/pd/kinds"
)

// [START Constructor Functions]

// newDateTime creates a dateTimeValue from atomic time.Time value
func newDateTime(val time.Time) dateTimeValue {
	if (time.Time{}) == val {
		return dateTimeValue{val, true}
	}
	return dateTimeValue{val, false}
}

// newSliceDateTime converts []time.Time{} -> Factory with dateTimeValues
func newSliceDateTime(vals []time.Time) Factory {
	var ret dateTimeValues
	for _, val := range vals {
		ret = append(ret, newDateTime(val))
	}
	return Factory{ret, kinds.DateTime}
}

// [END Constructor Functions]

// [START Converters]

// toFloat converts dateTimeValues to float64Values of the Unix EPOCH timestamp
// (seconds since midnight January 1, 1970)
// 2019-05-01 00:00:00 +0000 UTC: 1556757505
func (val dateTimeValue) toFloat64() float64Value {
	if val.null {
		return float64Value{math.NaN(), true}
	}
	v := val.v.UnixNano()
	return float64Value{float64(v), false}
}

// ToInt converts dateTimeValues to int64Values of the Unix EPOCH timestamp
// (seconds since midnight January 1, 1970)
//
// 2019-05-01 00:00:00 +0000 UTC: 1556757505
func (val dateTimeValue) toInt64() int64Value {
	if val.null {
		return int64Value{0, true}
	}
	v := val.v.UnixNano()
	return int64Value{v, false}
}

// ToBool converts dateTimeValues to boolValues
//
// x != time.Time{}: true; x == time.Time{}: false; null: false
func (val dateTimeValue) toBool() boolValue {
	if val.null {
		return boolValue{false, true}
	}
	if val.v == (time.Time{}) {
		return boolValue{false, false}
	}
	return boolValue{true, false}

}

// ToDateTime returns itself
func (val dateTimeValue) toDateTime() dateTimeValue {
	return val
}

// [END Converters]
