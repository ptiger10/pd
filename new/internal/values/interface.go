package values

import (
	"fmt"
	"math"
	"reflect"
	"time"

	"github.com/ptiger10/pd/new/options"
)

// [START Definitions]

// InterfaceValues is a slice of interface-typed Value/Null structs
type InterfaceValues []InterfaceValue

// An InterfaceValue is one interface-typed Value/Null struct
type InterfaceValue struct {
	V    interface{}
	Null bool
}

// Interface constructs an InterfaceValue
func Interface(v interface{}, null bool) InterfaceValue {
	return InterfaceValue{
		V:    v,
		Null: null,
	}
}

// [END Definitions]

// [START Converters]

// ToFloat converts InterfaceValues to FloatValues
// First checks for numerics, then tries to parse as a string
// If those fail, returns nil
//
func (vals InterfaceValues) ToFloat() Values {
	var ret FloatValues
	for _, val := range vals {
		if val.Null {
			ret = append(ret, Float(math.NaN(), true))
		} else {
			switch val.V.(type) {
			case float32, float64:
				v := reflect.ValueOf(val.V).Float()
				ret = append(ret, Float(v, false))
			case int, int8, int16, int32, int64:
				v := reflect.ValueOf(val.V).Int()
				ret = append(ret, Float(float64(v), false))
			case uint, uint8, uint16, uint32, uint64:
				v := reflect.ValueOf(val.V).Uint()
				ret = append(ret, Float(float64(v), false))
			default:
				vStr, ok := val.V.(string)
				if !ok {
					ret = append(ret, Float(math.NaN(), true))
					continue
				}
				ret = append(ret, stringToFloat(vStr))
			}
		}
	}
	return ret
}

// ToInt converts InterfaceValues to IntValues
// First checks for numerics, then tries to parse as a string
// If those fail, returns nil
//
func (vals InterfaceValues) ToInt() Values {
	var ret IntValues
	for _, val := range vals {
		if val.Null {
			ret = append(ret, Int(0, true))
		} else {
			switch val.V.(type) {
			case float32, float64:
				v := reflect.ValueOf(val.V).Float()
				ret = append(ret, Int(int64(v), false))
			case int, int8, int16, int32, int64:
				v := reflect.ValueOf(val.V).Int()
				ret = append(ret, Int(v, false))
			case uint, uint8, uint16, uint32, uint64:
				v := reflect.ValueOf(val.V).Uint()
				ret = append(ret, Int(int64(v), false))
			default:
				vStr, ok := val.V.(string)
				if !ok {
					ret = append(ret, Int(0, true))
					continue
				}
				ret = append(ret, stringToInt(vStr))
			}
		}
	}
	return ret
}

// ToBool converts InterfaceValues to BoolValues
//
// null: false; notnull: true
func (vals InterfaceValues) ToBool() Values {
	var ret BoolValues
	for _, val := range vals {
		if val.Null {
			ret = append(ret, Bool(false, true))
		} else {
			ret = append(ret, Bool(true, false))
		}
	}
	return ret
}

// ToDateTime converts InterfaceValues to DateTimeValues
//
// null: false; notnull: true
func (vals InterfaceValues) ToDateTime() Values {
	var ret DateTimeValues
	for _, val := range vals {
		if val.Null {
			ret = append(ret, DateTime(time.Time{}, true))
		} else {
			switch val.V.(type) {
			case float32, float64:
				f := reflect.ValueOf(val.V).Float()
				ret = append(ret, floatToDateTime(f))
			case int, int8, int16, int32, int64:
				i := reflect.ValueOf(val.V).Int()
				ret = append(ret, intToDateTime(i))
			case uint, uint8, uint16, uint32, uint64:
				u := reflect.ValueOf(val.V).Uint()
				ret = append(ret, intToDateTime(int64(u)))
			default:
				vStr, ok := val.V.(string)
				if !ok {
					ret = append(ret, DateTime(time.Time{}, true))
					continue
				}
				ret = append(ret, stringToDateTime(vStr))
			}
		}
	}
	return ret
}

// [END Converters]

// [START Methods]

// Describe the values in the collection
func (vals InterfaceValues) Describe() string {
	offset := options.DisplayValuesWhitespaceBuffer
	l := len(vals)
	len := fmt.Sprintf("%-*s%d\n", offset, "len", l)
	return fmt.Sprint(len)
}

// [END Methods]
