package values

import (
	"fmt"
	"math"
	"time"

	"github.com/ptiger10/pd/options"
)

// newFloat64 creates a float64Value from atomic float64 value
func newFloat64(val float64) float64Value {
	if math.IsNaN(val) {
		return float64Value{val, true}
	}
	return float64Value{val, false}
}

func (vals *float64Values) Less(i, j int) bool {
	if (*vals)[i].v < (*vals)[j].v {
		return true
	}
	return false
}

// [START Converters]

// toFloat returns itself
func (val float64Value) toFloat64() float64Value {
	return val
}

// toInt converts a float64Value to int64Value
//
// 1.9: 1, 1.5: 1, null: 0
func (val float64Value) toInt64() int64Value {
	if val.null {
		return int64Value{0, true}
	}
	v := int64(val.v)
	return int64Value{v, false}

}

func (val float64Value) toString() stringValue {
	if val.null {
		return stringValue{options.GetDisplayStringNullFiller(), true}
	}
	return stringValue{fmt.Sprint(val.v), false}
}

// toBool converts float64Value to boolValue
//
// x != 0: true; x == 0: false; null: false
func (val float64Value) toBool() boolValue {
	if val.null {
		return boolValue{false, true}
	}
	if val.v == 0 {
		return boolValue{false, false}
	}
	return boolValue{true, false}
}

// toDateTime converts float64Value to dateTimeValue.
// Tries to convert from Unix EPOCH time, otherwise returns null
func (val float64Value) toDateTime() dateTimeValue {
	if val.null {
		return dateTimeValue{time.Time{}, true}
	}
	return floatToDateTime(val.v)
}

func floatToDateTime(f float64) dateTimeValue {
	return intToDateTime(int64(f))
}

// [END Converters]
