package values

import (
	"math"
	"time"
)

// [START Constructor Functions]

// newBool creates a boolValue from atomic bool value
func newBool(val bool) boolValue {
	return boolValue{val, false}
}

func (vals *boolValues) Less(i, j int) bool {
	if (*vals)[i].v && !(*vals)[j].v {
		return true
	}
	return false
}

// [END Constructor Functions]

// [START Converters]
// toFloat converts boolValues to float64Values.
//
// true: 1.0, false: 0.0, null: NaN
func (val boolValue) toFloat64() float64Value {
	if val.null {
		return float64Value{math.NaN(), true}
	} else if val.v {
		return float64Value{1, false}
	} else {
		return float64Value{0, false}
	}
}

// toInt converts boolValues to int64Values.
//
// true: 1, false: 0, null: 0
func (val boolValue) toInt64() int64Value {
	if val.null {
		return int64Value{0, true}
	} else if val.v {
		return int64Value{1, false}
	} else {
		return int64Value{0, false}
	}
}

// toBool returns itself.
func (val boolValue) toBool() boolValue {
	return val
}

func (val boolValue) toDateTime() dateTimeValue {
	epochDate := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	if val.null {
		return dateTimeValue{time.Time{}, true}
	}
	return dateTimeValue{epochDate, false}
}

// [END Converters]
