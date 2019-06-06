package values

import (
	"math"
	"time"
)

// [START Constructor Functions]

// newInt64 creates an int64Value from atomic int64 value
func newInt64(val int64) int64Value {
	return int64Value{val, false}
}

// [END Constructor Functions]

// [START Converters]

// toFloat converts int64Value to float64Value
//
// 1: 1.0
func (val int64Value) toFloat64() float64Value {
	if val.null {
		return float64Value{math.NaN(), true}
	}
	v := float64(val.v)
	return float64Value{v, false}

}

// toInt returns itself
func (val int64Value) toInt64() int64Value {
	return val
}

// toBool converts int64Value to boolValue
//
// x != 0: true; x == 0: false; null: false
func (val int64Value) toBool() boolValue {
	if val.null {
		return boolValue{false, true}
	}
	if val.v == 0 {
		return boolValue{false, false}
	}
	return boolValue{true, false}
}

// toDateTime converts int64Value to dateTimeValue.
// Tries to convert from Unix EPOCH timestamp.
// Defaults to 1970-01-01 00:00:00 +0000 UTC.
func (val int64Value) toDateTime() dateTimeValue {
	if val.null {
		return dateTimeValue{time.Time{}, true}
	}
	return intToDateTime(val.v)
}

func intToDateTime(i int64) dateTimeValue {
	// convert from nanoseconds to seconds
	i /= 1000000000
	v := time.Unix(i, 0).UTC()
	return dateTimeValue{v, false}
}

// [END Converters]
