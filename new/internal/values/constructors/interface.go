package constructors

import (
	"math"
	"reflect"

	"github.com/ptiger10/pd/new/internal/values"
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

// SliceInterface converts []interface -> values.InterfaceValues
func SliceInterface(data interface{}) values.InterfaceValues {
	var vals values.InterfaceValues
	d := data.([]interface{})
	for i := 0; i < len(d); i++ {
		val := d[i]
		vals = append(vals)
		if isNullInterface(val) {
			vals = append(vals, values.Interface(val, true))
		} else {
			vals = append(vals, values.Interface(val, false))
		}

	}
	return vals
}
