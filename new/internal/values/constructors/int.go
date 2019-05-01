package constructors

import (
	"math"
	"reflect"
	"strconv"

	"github.com/ptiger10/pd/new/internal/values"
)

// [START Utilities]

// [END Utilities]

// [START Constructor Functions]

// SliceInt converts []int (of any variety) -> values.IntValues
func SliceInt(data interface{}) values.IntValues {
	var vals values.IntValues
	d := reflect.ValueOf(data)
	for i := 0; i < d.Len(); i++ {
		val := d.Index(i).Int()
		vals = append(vals, values.Int(val, false))
	}
	return vals
}

// SliceUInt converts []uint (of any variety) -> values.IntValues
func SliceUInt(data interface{}) values.IntValues {
	var vals values.IntValues
	d := reflect.ValueOf(data)
	for i := 0; i < d.Len(); i++ {
		val := d.Index(i).Uint()
		vals = append(vals, values.Int(int64(val), false))
	}
	return vals
}

// InterfaceInt converts []interface{} -> values.IntValues.
// Any values that cannot be converted into int64 are considered null
func InterfaceInt(d reflect.Value) values.IntValues {
	var vals values.IntValues
	for i := 0; i < d.Len(); i++ {
		v := d.Index(i).Elem()
		switch v.Kind() {
		case reflect.Invalid:
			vals = append(vals, values.Int(0, true))
		case reflect.String:
			val, err := strconv.ParseFloat(v.String(), 64)
			if err != nil || math.IsNaN(val) {
				vals = append(vals, values.Int(0, true))
			} else {
				vals = append(vals, values.Int(int64(val), false))
			}
		case reflect.Float32, reflect.Float64:
			val := v.Float()
			if math.IsNaN(val) {
				vals = append(vals, values.Int(0, true))
			} else {
				vals = append(vals, values.Int(int64(val), false))
			}

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			vals = append(vals, values.Int(int64(v.Int()), false))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			vals = append(vals, values.Int(int64(v.Uint()), false))
		default:
			vals = append(vals, values.Int(0, true))
		}
	}
	return vals
}

// [START Constructor Functions]
