package values

import (
	"math"
	"reflect"
	"time"

	"github.com/ptiger10/pd/kinds"
)

// [START Convenience Functions]

func isNullInterface(i interface{}) bool {
	switch i.(type) {
	case string:
		s := i.(string)
		if isNullString(s) {
			return true
		}
	case float32, float64:
		f := reflect.ValueOf(i).Float()
		if math.IsNaN(f) {
			return true
		}
	}
	return false
}

// [END Convenience Functions]

// newInterface creates an interfaceValue from atomic interface{} value
func newInterface(val interface{}) interfaceValue {
	if isNullInterface(val) {
		return interfaceValue{val, true}
	}
	return interfaceValue{val, false}
}

// newSliceInterface converts []interface -> Factory with interfaceValues
func newSliceInterface(vals []interface{}) Factory {
	var v interfaceValues
	for _, val := range vals {
		v = append(v, newInterface(val))
	}
	return Factory{v, kinds.Interface}
}

// [START Converters]
func (val interfaceValue) toFloat64() float64Value {
	if val.null {
		return float64Value{math.NaN(), true}
	}
	switch val.v.(type) {
	case float32, float64:
		v := reflect.ValueOf(val.v).Float()
		return newFloat(v)
	case int, int8, int16, int32, int64:
		v := reflect.ValueOf(val.v).Int()
		return newInt(v).toFloat64()
	case uint, uint8, uint16, uint32, uint64:
		v := reflect.ValueOf(val.v).Uint()
		return newInt(int64(v)).toFloat64()
	case string:
		return newString(val.v.(string)).toFloat64()
	case bool:
		return newBool(val.v.(bool)).toFloat64()
	case time.Time:
		return newDateTime(val.v.(time.Time)).toFloat64()
	}
	return float64Value{}
}

func (val interfaceValue) toInt64() int64Value {
	if val.null {
		return int64Value{0, true}
	}
	switch val.v.(type) {
	case float32, float64:
		v := reflect.ValueOf(val.v).Float()
		return newFloat(v).toInt64()
	case int, int8, int16, int32, int64:
		v := reflect.ValueOf(val.v).Int()
		return newInt(v)
	case uint, uint8, uint16, uint32, uint64:
		v := reflect.ValueOf(val.v).Uint()
		return int64Value{int64(v), false}
	case string:
		return newString(val.v.(string)).toInt64()
	case bool:
		return newBool(val.v.(bool)).toInt64()
	case time.Time:
		return newDateTime(val.v.(time.Time)).toInt64()
	}
	return int64Value{}
}

func (val interfaceValue) toBool() boolValue {
	if val.null {
		return boolValue{false, true}
	}
	switch val.v.(type) {
	case float32, float64:
		v := reflect.ValueOf(val.v).Float()
		return newFloat(v).toBool()
	case int, int8, int16, int32, int64:
		v := reflect.ValueOf(val.v).Int()
		return newInt(v).toBool()
	case uint, uint8, uint16, uint32, uint64:
		v := reflect.ValueOf(val.v).Uint()
		return newInt(int64(v)).toBool()
	case string:
		return newString(val.v.(string)).toBool()
	case bool:
		return newBool(val.v.(bool))
	case time.Time:
		return newDateTime(val.v.(time.Time)).toBool()
	}
	return boolValue{}
}

func (val interfaceValue) toDateTime() dateTimeValue {
	if val.null {
		return dateTimeValue{time.Time{}, true}
	}
	switch val.v.(type) {
	case float32, float64:
		v := reflect.ValueOf(val.v).Float()
		return newFloat(v).toDateTime()
	case int, int8, int16, int32, int64:
		v := reflect.ValueOf(val.v).Int()
		return newInt(v).toDateTime()
	case uint, uint8, uint16, uint32, uint64:
		v := reflect.ValueOf(val.v).Uint()
		return newInt(int64(v)).toDateTime()
	case string:
		return newString(val.v.(string)).toDateTime()
	case bool:
		return newBool(val.v.(bool)).toDateTime()
	case time.Time:
		return newDateTime(val.v.(time.Time))
	}
	return dateTimeValue{}
}

func (val interfaceValue) toInterface() interfaceValue {
	return val
}

// [END Converters]
