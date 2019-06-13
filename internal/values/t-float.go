package values

import (
	"math"
	"time"
)

// [START Constructor Functions]

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

// func (vals *float64Values) Filter(fn func(interface{}) bool) Values {
// 	var newVals float64Values
// 	for _, val := range *vals {
// 		if fn() {

// 		}
// 	}
// }

// func (vals dateTimeValues) Filter(callbackFn func(time.Time) bool) dateTimeValues {
// 	var ret dateTimeValues
// 	valid, _ := vals.valid()
// 	for _, val := range valid {
// 		if callbackFn(val) {
// 			ret = append(ret, dateTimeValue{v: val})
// 		}
// 	}
// 	return ret
// }

// func (s *Series) FilterDateTime(callbackFn func(time.Time) bool) (*Series, error) {
// 	if s.Kind != DateTime {
// 		return s, fmt.Errorf("FilterString can be called only on Series with type String, not %v", s.Kind)
// 	}
// 	vals := s.Values.(dateTimeValues).Filter(callbackFn)

// 	return Series{
// 		Values: vals,
// 		DataType:   DateTime,
// 	}, nil
// }

// [END Constructor Functions]

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
