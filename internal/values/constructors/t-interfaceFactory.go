package constructors

import (
	"math"
	"reflect"

	"github.com/ptiger10/pd/internal/values"
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

// SliceInterface converts []interface -> ValuesFactory with values.InterfaceValues
func SliceInterface(vals []interface{}) ValuesFactory {
	var v values.InterfaceValues
	for _, val := range vals {
		if isNullInterface(val) {
			v = append(v, values.Interface(val, true))
		} else {
			v = append(v, values.Interface(val, false))
		}

	}
	return ValuesFactory{v, kinds.Interface}
}
