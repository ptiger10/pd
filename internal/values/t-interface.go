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

// SliceInterface converts []interface -> Factory with interfaceValues
func SliceInterface(vals []interface{}) Factory {
	var v interfaceValues
	for _, val := range vals {
		if isNullInterface(val) {
			v = append(v, interfaceVal(val, true))
		} else {
			v = append(v, interfaceVal(val, false))
		}

	}
	return Factory{v, kinds.Interface}
}

// [START Converters]

// ToFloat converts interfaceValues to float64Values
// First checks for numerics, then tries to parse as a string
// If those fail, returns nil
//
func (vals interfaceValues) ToFloat() Values {
	var ret float64Values
	for _, val := range vals {
		if val.null {
			ret = append(ret, float64Val(math.NaN(), true))
		} else {
			switch val.v.(type) {
			case float32, float64:
				v := reflect.ValueOf(val.v).Float()
				ret = append(ret, float64Val(v, false))
			case int, int8, int16, int32, int64:
				v := reflect.ValueOf(val.v).Int()
				ret = append(ret, float64Val(float64(v), false))
			case uint, uint8, uint16, uint32, uint64:
				v := reflect.ValueOf(val.v).Uint()
				ret = append(ret, float64Val(float64(v), false))
			default:
				vStr, ok := val.v.(string)
				if !ok {
					ret = append(ret, float64Val(math.NaN(), true))
					continue
				}
				ret = append(ret, stringToFloat(vStr))
			}
		}
	}
	return ret
}

// ToInt converts interfaceValues to int64Values
// First checks for numerics, then tries to parse as a string
// If those fail, returns nil
//
func (vals interfaceValues) ToInt() Values {
	var ret int64Values
	for _, val := range vals {
		if val.null {
			ret = append(ret, int64Val(0, true))
		} else {
			switch val.v.(type) {
			case float32, float64:
				v := reflect.ValueOf(val.v).Float()
				ret = append(ret, int64Val(int64(v), false))
			case int, int8, int16, int32, int64:
				v := reflect.ValueOf(val.v).Int()
				ret = append(ret, int64Val(v, false))
			case uint, uint8, uint16, uint32, uint64:
				v := reflect.ValueOf(val.v).Uint()
				ret = append(ret, int64Val(int64(v), false))
			default:
				vStr, ok := val.v.(string)
				if !ok {
					ret = append(ret, int64Val(0, true))
					continue
				}
				ret = append(ret, stringToInt(vStr))
			}
		}
	}
	return ret
}

// ToBool converts interfaceValues to boolValues
//
// null: false; notnull: true
func (vals interfaceValues) ToBool() Values {
	var ret boolValues
	for _, val := range vals {
		if val.null {
			ret = append(ret, boolVal(false, true))
		} else {
			switch val.v.(type) {
			case bool:
				v := val.v.(bool)
				ret = append(ret, boolVal(v, false))
			default:
				ret = append(ret, boolVal(false, true))
			}

		}
	}
	return ret
}

// ToDateTime converts interfaceValues to dateTimeValues
//
// null: false; notnull: true
func (vals interfaceValues) ToDateTime() Values {
	var ret dateTimeValues
	for _, val := range vals {
		if val.null {
			ret = append(ret, dateTimeVal(time.Time{}, true))
		} else {
			switch val.v.(type) {
			case float32, float64:
				f := reflect.ValueOf(val.v).Float()
				ret = append(ret, floatToDateTime(f))
			case int, int8, int16, int32, int64:
				i := reflect.ValueOf(val.v).Int()
				ret = append(ret, intToDateTime(i))
			case uint, uint8, uint16, uint32, uint64:
				u := reflect.ValueOf(val.v).Uint()
				ret = append(ret, intToDateTime(int64(u)))
			case time.Time:
				t := val.v.(time.Time)
				if t == (time.Time{}) {
					ret = append(ret, dateTimeVal(time.Time{}, true))
				} else {
					ret = append(ret, dateTimeVal(t, false))
				}

			default:
				vStr, ok := val.v.(string)
				if !ok {
					ret = append(ret, dateTimeVal(time.Time{}, true))
					continue
				}
				ret = append(ret, stringToDateTime(vStr))
			}
		}
	}
	return ret
}

// [END Converters]
